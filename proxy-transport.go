package momos

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type proxyTransport struct {
	http.RoundTripper
}

// RoundTrip will replace "server" with "schmerver"
func (t *proxyTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	// get the response of an given request
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	// Just an example

	element := doc.Find("ssi-teaser")

	eHTML, err := element.Html()

	if err != nil {
		panic(err)
	}

	element.SetHtml(strings.Replace("schmerver", "server", eHTML, -1))

	htmlDoc, err := doc.Html()

	if err != nil {
		panic(err)
	}

	// assign new reader with content
	content := []byte(htmlDoc)
	body := ioutil.NopCloser(bytes.NewReader(content))
	resp.Body = body
	resp.ContentLength = int64(len(content)) // update content length
	resp.Header.Set("Content-Length", strconv.Itoa(len(content)))
	return resp, nil
}
