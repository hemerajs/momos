package momos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		se := SSIElement{Element: element}
		se.Pos = element.Index()
		se.Attributes = SSIAttributes{
			"timeout":  element.AttrOr("timeout", "2000"),
			"src":      element.AttrOr("src", ""),
			"fallback": element.AttrOr("fallback", ""),
			"cache":    element.AttrOr("cache", ""),
			"name":     element.AttrOr("name", ""),
		}

		se.GetError()
		se.GetTimeout()

		// get SSI content and replace tag
		if fragmentUrl, ok := se.Attributes["src"]; ok {
			// @TODO parse timeout
			timeout := time.Duration(2000 * time.Millisecond)
			fmt.Printf("Call: %v\n", fragmentUrl)
			client := http.Client{Timeout: timeout}
			resp, err := client.Get(fragmentUrl)

			//@TODO error handling
			if err != nil {
				html, err := se.Error.Html()

				//@TODO error handling
				if err == nil {
					element.ReplaceWithHtml(html)
				}
			} else {
				content, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					panic(err)
				}

				element.ReplaceWithHtml(string(content))
				resp.Body.Close()
			}
		}

		eHTML, err := element.Html()

		//@TODO error handling
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

	// assign new reader with content
	content := []byte(htmlDoc)
	body := ioutil.NopCloser(bytes.NewReader(content))
	resp.Body = body
	resp.ContentLength = int64(len(content)) // update content length
	resp.Header.Set("Content-Length", strconv.Itoa(len(content)))
	return resp, nil
}
