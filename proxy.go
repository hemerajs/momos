package momos

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	reverseProxy *httputil.ReverseProxy
	proxyUrl     string
}

func New(proxyUrl, targetUrl string) *Proxy {
	target, tErr := url.Parse(targetUrl)

	if tErr != nil {
		panic(tErr)
	}

	p := &Proxy{}
	p.proxyUrl = proxyUrl
	p.reverseProxy = httputil.NewSingleHostReverseProxy(target)
	p.reverseProxy.Transport = &proxyTransport{http.DefaultTransport}

	return p
}

func (p *Proxy) Start() error {
	return http.ListenAndServe(p.proxyUrl, p.reverseProxy)
}
