package momos

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	ssiErrorTag    = "ssi-error"
	ssiTimeoutTag  = "ssi-timeout"
	defaultTimeout = 2000
)

type SSIAttributes map[string]string

type SSIElement struct {
	Tag           string
	HasErrorTag   bool
	errorTag      *goquery.Selection
	HasTimeoutTag bool
	timeoutTag    *goquery.Selection
	Attributes    SSIAttributes
	Element       *goquery.Selection
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

func (s *SSIElement) replaceWithDefaultHTML() error {
	s.Element.Find(ssiErrorTag + "," + ssiTimeoutTag).Remove()

	html, err := s.Element.Html()

	if err == nil {
		s.Element.ReplaceWithHtml(html)
	} else {
		return err
	}

	return nil
}

func (s *SSIElement) replaceWithErrorHTML() error {
	if s.HasErrorTag {
		html, err := s.errorTag.Html()

		if err == nil {
			s.Element.ReplaceWithHtml(html)
		} else {
			return err
		}
	} else {
		s.replaceWithDefaultHTML()
	}

	return nil
}

func (s *SSIElement) replaceWithTimeoutHTML() error {
	if s.HasTimeoutTag {
		html, err := s.timeoutTag.Html()

		if err == nil {
			s.Element.ReplaceWithHtml(html)
		} else {
			return err
		}
	} else {
		s.replaceWithDefaultHTML()
	}

	return nil
}

func (s *SSIElement) Timeout() (int, error) {
	timeoutMs, err := strconv.Atoi(s.Attributes["timeout"])

	if err != nil {
		errorf("illegal value %q in timeout attribute", timeoutMs)
		return timeoutMs, nil
	}

	return defaultTimeout, err
}

func (s *SSIElement) SetupSuccess(body []byte) {
	s.Element.ReplaceWithHtml(string(body))
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
