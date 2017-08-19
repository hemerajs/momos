package momos

import (
	"bytes"
	"html/template"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	ssiErrorTag    = "ssi-error"
	ssiTimeoutTag  = "ssi-timeout"
	defaultTimeout = 2000
)

type TemplateContext struct {
	DateLocal string
	Date      string
	RequestId string
}

type SSIAttributes map[string]string

type SSIElement struct {
	Tag             string
	HasErrorTag     bool
	errorTag        *goquery.Selection
	HasTimeoutTag   bool
	timeoutTag      *goquery.Selection
	Attributes      SSIAttributes
	Element         *goquery.Selection
	templateContext TemplateContext
	hasTemplate     bool
}

func (s *SSIElement) GetErrorTag() error {
	node := s.Element.Find(ssiErrorTag)
	s.HasErrorTag = node.Length() > 0
	s.timeoutTag = node

	return nil
}

func (s *SSIElement) GetTimeoutTag() error {
	node := s.Element.Find(ssiTimeoutTag)
	s.HasTimeoutTag = node.Length() > 0
	s.errorTag = node

	return nil
}

func (s *SSIElement) Url() string {
	return s.Attributes["src"]
}

func (s *SSIElement) Name() string {
	return s.Attributes["name"]
}

func (s *SSIElement) HasHemplate() bool {
	return s.Attributes["template"] == "true"
}

func (s *SSIElement) Timeout() (int, error) {
	timeoutMs, err := strconv.Atoi(s.Attributes["timeout"])

	if err != nil {
		errorf("illegal value %q in timeout attribute", timeoutMs)
		return timeoutMs, nil
	}

	return defaultTimeout, err
}

func (s *SSIElement) replaceWithDefaultHTML() error {
	s.Element.Find(ssiErrorTag + "," + ssiTimeoutTag).Remove()

	html, err := s.Element.Html()

	if err == nil {
		err := s.ReplaceWithHtml(html)
		if err != nil {
			s.replaceWithDefaultHTML()
			return err
		}
	} else {
		return err
	}

	return nil
}

func (s *SSIElement) replaceWithErrorHTML() error {
	if s.HasErrorTag {
		html, err := s.errorTag.Html()

		if err == nil {
			err := s.ReplaceWithHtml(html)
			if err != nil {
				s.replaceWithDefaultHTML()
				return err
			}
		} else {
			return err
		}
	} else {
		s.replaceWithDefaultHTML()
	}

	return nil
}

func (s *SSIElement) ReplaceWithHtml(html string) error {
	if s.HasHemplate() {
		var doc bytes.Buffer
		tpl, err := template.New(s.Name()).Parse(html)

		if err != nil {
			errorf("template parsing error %q", err)
			return err
		}

		err = tpl.Execute(&doc, s.templateContext)

		if err != nil {
			errorf("Error during template rendering %q", err)
			return err
		}

		s.Element.ReplaceWithHtml(doc.String())
	} else {
		s.Element.ReplaceWithHtml(html)
	}

	return nil
}

func (s *SSIElement) replaceWithTimeoutHTML() error {
	if s.HasTimeoutTag {
		html, err := s.timeoutTag.Html()

		if err == nil {
			err := s.ReplaceWithHtml(html)
			if err != nil {
				s.replaceWithDefaultHTML()
				return err
			}
		} else {
			return err
		}
	} else {
		s.replaceWithDefaultHTML()
	}

	return nil
}

func (s *SSIElement) SetupSuccess(body []byte) error {
	return s.ReplaceWithHtml(string(body))
}

func (s *SSIElement) SetupFallback(err error) error {
	switch err {
	case ErrInvalidContentType:
	case ErrInvalidStatusCode:
		s.replaceWithErrorHTML()
	case ErrRequest:
		s.replaceWithErrorHTML()
	case ErrTimeout:
		s.replaceWithTimeoutHTML()
	}

	return err
}
