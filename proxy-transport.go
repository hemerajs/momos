package momos

import (
	"bytes"
	"fmt"
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

	doc.Find("ssi").Each(func(i int, element *goquery.Selection) {
		se := SSIElement{}
		se.Pos = element.Index()
		se.Attributes = SSIAttributes{
			"timeout":  element.AttrOr("timeout", "2000"),
			"url":      element.AttrOr("timeout", ""),
			"fallback": element.AttrOr("fallback", ""),
			"cache":    element.AttrOr("cache", ""),
			"name":     element.AttrOr("name", ""),
		}

		elementErrorTag := element.Find("ssi-error")
		se.HasErrorTag = elementErrorTag.Length() > 0

		elementErrorTagHTML, tagErr := elementErrorTag.Html()

		if tagErr != nil {
			panic(err)
		}

		se.ErrorHTML = elementErrorTagHTML

		eHTML, err := element.Html()

		if err != nil {
			panic(err)
		}

		se.Len = len(eHTML)

		element.SetHtml(strings.Replace(eHTML, "server", "schmerver", -1))
	})

	htmlDoc, err := doc.Html()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", htmlDoc)

	// assign new reader with content
	content := []byte(htmlDoc)
	body := ioutil.NopCloser(bytes.NewReader(content))
	resp.Body = body
	resp.ContentLength = int64(len(content)) // update content length
	resp.Header.Set("Content-Length", strconv.Itoa(len(content)))
	return resp, nil
}
