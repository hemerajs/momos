package momos

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type proxyTransport struct {
	http.RoundTripper
}

// RoundTrip will replace "server" with "schmerver"
func (t *proxyTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	timeStart := time.Now()

	// get the response of an given request
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		errorf("could not create RoundTripper from %q", req.URL)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		errorf("illegal response body from %q", req.URL)
		return nil, err
	}

	doc.Find("ssi").Each(func(i int, element *goquery.Selection) {
		se := SSIElement{Element: element}
		se.Attributes = SSIAttributes{
			"timeout":  element.AttrOr("timeout", "2000"),
			"src":      element.AttrOr("src", ""),
			"fallback": element.AttrOr("fallback", ""),
			"cache":    element.AttrOr("cache", ""),
			"name":     element.AttrOr("name", ""),
		}

		se.GetErrorHTML()
		se.GetTimeoutHTML()

		err := se.makeRequest()

		if err != nil {
			errorf("ssi error for url: %q, name: %q", req.URL, se.Attributes["name"])
		}
	})

	htmlDoc, err := doc.Html()

	if err != nil {
		errorf("Could not get html from document %q", req.URL)
		return nil, err
	}

	// assign new reader with content
	content := []byte(htmlDoc)
	body := ioutil.NopCloser(bytes.NewReader(content))
	resp.Body = body
	resp.ContentLength = int64(len(content)) // update content length
	resp.Header.Set("Content-Length", strconv.Itoa(len(content)))

	debugf("Process Complete Request %q took %q", req.URL, time.Since(timeStart))

	return resp, nil
}
