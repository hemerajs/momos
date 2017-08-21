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

	"github.com/nats-io/nuid"

	"github.com/PuerkitoBio/goquery"
	clientCache "github.com/gregjones/httpcache"
)

var (
	ErrRequest            = errors.New("Request error")
	ErrTimeout            = errors.New("Timeout error")
	ErrInvalidStatusCode  = errors.New("Invalid status code")
	ErrInvalidContentType = errors.New("Invalid content type")
)

var cache = clientCache.NewMemoryCacheTransport()

type ssiResult struct {
	name    string
	payload []byte
	error   error
}

type proxyTransport struct {
	http.RoundTripper
}

func (t *proxyTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	timeStart := time.Now()

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		Log.Errorf("could not create RoundTripper from %q", req.URL)
		return nil, err
	}

	// Only html files are scanned
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return resp, nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		Log.Errorf("illegal response body from %q", req.URL)
		return nil, err
	}

	ssiCount := 0
	ssiElements := map[string]SSIElement{}

	doc.Find("ssi").Each(func(i int, element *goquery.Selection) {
		se := SSIElement{Element: element}

		se.SetTimeout(element.AttrOr("timeout", "2000"))
		se.SetSrc(element.AttrOr("src", ""))
		se.SetName(element.AttrOr("name", nuid.Next()))
		se.SetTemplate(element.AttrOr("template", "false"))
		se.SetFilterIncludes(element.AttrOr("no-scripts", "true"))

		se.GetErrorTag()
		se.GetTimeoutTag()

		se.templateContext = TemplateContext{
			DateLocal: time.Now().Local().Format("2006-01-02"),
			Date:      time.Now().Format(time.RFC3339),
			RequestId: req.Header.Get("X-Request-Id"),
			Name:      se.name,
		}

		ssiElements[se.name] = se
		ssiCount++
	})

	ch := make(chan ssiResult)

	timeStartRequest := time.Now()

	for _, el := range ssiElements {
		go makeRequest(el.name, el.src, ch, el.timeout)
	}

	for i := 0; i < ssiCount; i++ {
		res := <-ch
		el := ssiElements[res.name]
		if res.error == nil {
			Log.Tracef("Call fragment %q, url: %v, duration: %v", el.name, el.src, time.Since(timeStartRequest))
			el.SetupSuccess(res.payload)
		} else {
			el.SetupFallback(res.error)
			Log.Errorf("Fragment error %q, url: %v\nerror: %q", el.name, el.src, res.error)
		}
	}

	close(ch)

	htmlDoc, err := doc.Html()

	if err != nil {
		Log.Errorf("could not get html from document %q", req.URL)
		return nil, err
	}

	// assign new reader with content
	content := []byte(htmlDoc)
	body := ioutil.NopCloser(bytes.NewReader(content))
	resp.Body = body
	resp.ContentLength = int64(len(content)) // update content length
	resp.Header.Set("Content-Length", strconv.Itoa(len(content)))

	Log.Tracef("Processing complete %q duration: %v", req.URL, time.Since(timeStart))

	return resp, nil
}

func makeRequest(name string, url string, ch chan<- ssiResult, timeout time.Duration) {
	// @TODO don't create a new client per request
	var Client = &http.Client{
		Transport: cache,
		Timeout:   timeout,
	}

	resp, err := Client.Get(url)

	if err, ok := err.(net.Error); ok && err.Timeout() {
		ch <- ssiResult{name: name, error: ErrTimeout}
	} else if err != nil {
		ch <- ssiResult{name: name, error: err}
	} else {
		contentType := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "text/html") {
			ch <- ssiResult{name: name, error: ErrInvalidContentType}
		} else if resp.StatusCode > 199 && resp.StatusCode < 300 {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			// https://github.com/gregjones/httpcache
			if resp.Header.Get("X-From-Cache") == "1" {
				Log.Noticef(`Fragment "%v" was cached`, name)
			} else {
				Log.Noticef(`Fragment "%v" was not cached`, name)
			}

			ch <- ssiResult{name: name, payload: body}
		} else {
			ch <- ssiResult{name: name, error: ErrInvalidStatusCode}
		}
	}
}
