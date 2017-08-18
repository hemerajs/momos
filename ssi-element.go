package momos

import "github.com/PuerkitoBio/goquery"

type SSIAttributes map[string]string

type SSIElement struct {
	Tag           string
	Pos           int
	Len           int
	HasErrorTag   bool
	Error         *goquery.Selection
	HasTimeoutTag bool
	Timeout       *goquery.Selection
	Attributes    SSIAttributes
	Element       *goquery.Selection
}

func (s *SSIElement) GetError() error {
	node := s.Element.Find("ssi-error")
	s.HasErrorTag = node.Length() > 0
	s.Timeout = node

	return nil
}

func (s *SSIElement) GetTimeout() error {
	node := s.Element.Find("ssi-error")
	s.HasTimeoutTag = node.Length() > 0
	s.Error = node

	return nil
}
