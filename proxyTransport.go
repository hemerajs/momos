package momos

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
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

	// read response buffer
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// close reader
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// replace "server" with "schmerver"
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	// assign new reader with content
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b)) // update content length
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return resp, nil
}
