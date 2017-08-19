package momos

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	clientCache "github.com/gregjones/httpcache"
)

var (
	ErrRequest            = errors.New("Request error")
	ErrTimeout            = errors.New("Timeout error")
	ErrInvalidStatusCode  = errors.New("Invalid status code")
	ErrInvalidContentType = errors.New("Invalid content type")
)

var Client = &http.Client{Transport: clientCache.NewMemoryCacheTransport()}

type proxyTransport struct {
	http.RoundTripper
}

func (t *proxyTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	timeStart := time.Now()

	// get the response of an given request
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		errorf("could not create RoundTripper from %q", req.URL)
		return nil, err
	}

	// Only html files are scanned
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return resp, nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		errorf("illegal response body from %q", req.URL)
		return nil, err
	}

	ssiElements := []SSIElement{}

	doc.Find("ssi").Each(func(i int, element *goquery.Selection) {
		se := SSIElement{Element: element}
		se.Attributes = SSIAttributes{
			"timeout": element.AttrOr("timeout", "2000"),
			"src":     element.AttrOr("src", ""),
			"name":    element.AttrOr("name", ""),
		}

		se.GetErrorTag()
		se.GetTimeoutTag()

		ssiElements = append(ssiElements, se)
	})

	ch := make(chan []byte)
	chErr := make(chan error)

	timeStartRequest := time.Now()

	for _, element := range ssiElements {
		timeout, _ := element.Timeout()
		go makeRequest(element.Url(), ch, chErr, timeout)
	}

	for _, element := range ssiElements {
		select {
		case res := <-ch:
			debugf("SSI [%v] - Request to %v took %v", element.Name(), element.Url(), time.Since(timeStartRequest))
			element.SetupSuccess(res)
		case err := <-chErr:
			element.SetupFallback(err)
			debugf("SSI [%v] - Request to %v error: %q", element.Name(), element.Url(), err)
		}
	}

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

func makeRequest(url string, ch chan<- []byte, chErr chan<- error, timeoutMs int) {
	timeout := time.Duration(time.Duration(timeoutMs) * time.Millisecond)
	Client.Timeout = timeout
	resp, err := Client.Get(url)

	if err != nil {
		chErr <- ErrRequest
	} else if err, ok := err.(net.Error); ok && err.Timeout() {
		chErr <- ErrTimeout
	} else {
		contentType := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "text/html") {
			chErr <- ErrInvalidContentType
		} else if resp.StatusCode > 199 && resp.StatusCode < 300 {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			ch <- body
		}
	}

}
