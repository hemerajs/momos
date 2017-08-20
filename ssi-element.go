package momos

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	ssiErrorTag         = "ssi-error"
	ssiTimeoutTag       = "ssi-timeout"
	jsIncludeSelection  = "script"
	cssIncludeSelection = "link"
	defaultTimeout      = 2000
)

type TemplateContext struct {
	DateLocal string
	Date      string
	RequestId string
	Name      string
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
	filterIncludes  bool
	name            string
	src             string
	timeout         time.Duration
	hasTemplate     bool
}

// GetErrorTag find the error tag
func (s *SSIElement) GetErrorTag() error {
	node := s.Element.Find(ssiErrorTag)
	s.HasErrorTag = node.Length() > 0
	s.errorTag = node

	return nil
}

// GetTimeoutTag find the timeout tag
func (s *SSIElement) GetTimeoutTag() error {
	node := s.Element.Find(ssiTimeoutTag)
	s.HasTimeoutTag = node.Length() > 0
	s.timeoutTag = node

	return nil
}

// SetName set the name of the ssi fragment
func (s *SSIElement) SetName(name string) error {
	s.name = name
	return nil
}

// SetSrc set the source of the ssi fragment
func (s *SSIElement) SetSrc(src string) error {
	s.src = src
	return nil
}

// SetTemplate enables go templating
func (s *SSIElement) SetTemplate(h string) error {
	if h == "false" {
		s.hasTemplate = false
	} else {
		s.hasTemplate = true
	}
	return nil
}

// SetFilterIncludes enable filter of dynamic includes like <script> and <link> tags
func (s *SSIElement) SetFilterIncludes(h string) error {
	if h == "false" {
		s.filterIncludes = false
	} else {
		s.filterIncludes = true
	}
	return nil
}

// SetTimeout set the timeout
func (s *SSIElement) SetTimeout(t string) error {
	timeoutMs, err := strconv.Atoi(t)

	if err != nil {
		errorf("illegal value %q in timeout attribute", timeoutMs)
		return err
	}

	s.timeout = time.Duration(time.Duration(timeoutMs) * time.Millisecond)

	return nil
}

// replaceWithDefaultHTML replace the fragment with the default content
// empty string (default)
func (s *SSIElement) replaceWithDefaultHTML() error {
	s.Element.Find(ssiErrorTag + "," + ssiTimeoutTag).Remove()

	html, err := s.Element.Html()

	if err == nil {
		return s.ReplaceWithHTML(html)
	}

	return nil
}

// removeDynamicIncludes removes all scripts and link includes
func removeDynamicIncludes(html string) (string, error) {

	r := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		errorf("could not parse as html document")
		return "", err
	}

	doc.Find(jsIncludeSelection + "," + cssIncludeSelection).Remove()

	html, err = doc.Html()

	if err != nil {
		errorf("could not get html from document")
		return "", err
	}

	return html, nil
}

// replaceWithErrorHTML use the content from the ssi-error tag as fallback
// if ssi-error could not be replaced we use the default content as fallback
func (s *SSIElement) replaceWithErrorHTML() error {
	if s.HasErrorTag {
		html, err := s.errorTag.Html()
		if err == nil {
			err := s.ReplaceWithHTML(html)
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

// ReplaceWithHTML use the content from the ssi service
// and parse the fragment as go template (optionally)
func (s *SSIElement) ReplaceWithHTML(html string) error {
	if s.hasTemplate {
		var doc bytes.Buffer
		tpl, err := template.New(s.name).Parse(html)

		if err != nil {
			errorf("template parsing error %q", err)
			return err
		}

		err = tpl.Execute(&doc, s.templateContext)

		if err != nil {
			errorf("error during template rendering %q", err)
			return err
		}

		s.Element.ReplaceWithHtml(doc.String())
	} else {
		s.Element.ReplaceWithHtml(html)
	}

	return nil
}

// replaceWithTimeoutHTML use the content from the ssi-timeout tag as fallback
// if ssi-timeout could not be replaced we use the default content as fallback
func (s *SSIElement) replaceWithTimeoutHTML() error {
	if s.HasTimeoutTag {
		html, err := s.timeoutTag.Html()

		if err == nil {
			err := s.ReplaceWithHTML(html)
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

// SetupSuccess replace the fragment with the correct content
func (s *SSIElement) SetupSuccess(body []byte) error {
	html := string(body)

	if s.filterIncludes {
		h, err := removeDynamicIncludes(html)

		if err != nil {
			errorf("could not remove dynamic includes from document %q", err)
			return s.replaceWithErrorHTML()
		}
		return s.ReplaceWithHTML(h)
	}

	return s.ReplaceWithHTML(html)
}

// SetupFallback replace the fragment with the correct fallback content
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
