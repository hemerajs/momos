package momos

import (
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	ssiErrorTag   = "ssi-error"
	ssiTimeoutTag = "ssi-timeout"
)

type SSIAttributes map[string]string

type SSIElement struct {
	Tag           string
	HasErrorTag   bool
	Error         *goquery.Selection
	HasTimeoutTag bool
	Timeout       *goquery.Selection
	Attributes    SSIAttributes
	Element       *goquery.Selection
}

func (s *SSIElement) GetErrorHTML() error {
	node := s.Element.Find(ssiErrorTag)
	s.HasErrorTag = node.Length() > 0
	s.Timeout = node

	return nil
}

func (s *SSIElement) GetTimeoutHTML() error {
	node := s.Element.Find(ssiTimeoutTag)
	s.HasTimeoutTag = node.Length() > 0
	s.Error = node

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
		html, err := s.Error.Html()

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
		html, err := s.Timeout.Html()

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

// makeRequest start a GET request to get the SSI content
// Any none 2XX status code is handled as an error
// Timeout errors are handled with the `ssi-timeout` tag
// Any other error is handled with the `ssi-error` tag
func (s *SSIElement) makeRequest() error {
	timeoutMs, err := strconv.Atoi(s.Attributes["timeout"])

	if err != nil {
		errorf("illegal value %q in timeout attribute", timeoutMs)
		return err
	}

	if fragmentURL, ok := s.Attributes["src"]; ok {
		timeout := time.Duration(time.Duration(timeoutMs) * time.Millisecond)
		client := http.Client{Timeout: timeout}
		timeStart := time.Now()
		resp, err := client.Get(fragmentURL)

		debugf("[%q] Request to %q took %q", s.Attributes["name"], fragmentURL, time.Since(timeStart))

		// Only html
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			errorf("ssi: invalid content type %q", fragmentURL, contentType)
			return nil
		}

		if err != nil { // request error
			errorf("Request error %q", fragmentURL)
			s.replaceWithErrorHTML()
		} else if err, ok := err.(net.Error); ok && err.Timeout() { // Timeout
			errorf("Request timeouts %q, timeout: %q", fragmentURL, timeout)
			s.replaceWithTimeoutHTML()
		} else if resp.StatusCode > 199 && resp.StatusCode < 300 { // None 2xx status code
			content, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				errorf("Could not read response from %q", fragmentURL)
				return err
			}
			// replace with content from ssi service
			s.Element.ReplaceWithHtml(string(content))
		} else { // Default
			errorf("invalid request code %v from %q", resp.StatusCode, fragmentURL)
			s.replaceWithErrorHTML()
		}
	}

	return nil
}
