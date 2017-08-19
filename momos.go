package momos

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/lox/httpcache"
	"github.com/lox/httpcache/httplog"
)

var ServerLogging = false

type Proxy struct {
	reverseProxy *httputil.ReverseProxy
	proxyURL     string
	Server       *http.Server
	handler      *httpcache.Handler
	cache        httpcache.Cache
}

func PreCacheResponseHandler(h http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		debugf("▨ PreResponse (%v) - Cache is %v", req.Host+req.URL.String(), res.Header().Get("X-Cache"))
		h.ServeHTTP(res, req)
	}
}

func New(proxyUrl, targetUrl string) *Proxy {
	target, tErr := url.Parse(targetUrl)

	if tErr != nil {
		errorf("Invalid url: %v", targetUrl)
		panic(tErr)
	}

	httpcache.DebugLogging = ServerLogging

	p := &Proxy{}
	p.proxyURL = proxyUrl
	p.reverseProxy = httputil.NewSingleHostReverseProxy(target)
	p.reverseProxy.Transport = &proxyTransport{http.DefaultTransport}

	p.cache = httpcache.NewMemoryCache()
	p.handler = httpcache.NewHandler(p.cache, PreCacheResponseHandler(p.reverseProxy))
	p.handler.Shared = true

	respLogger := httplog.NewResponseLogger(p.handler)
	respLogger.DumpRequests = true
	respLogger.DumpResponses = true
	respLogger.DumpErrors = true

	return p
}

// Start starts the server and listen on the given port
func (p *Proxy) Start() error {
	p.Server = &http.Server{Addr: p.proxyURL}
	p.Server.Handler = p.handler
	return p.Server.ListenAndServe()
}

// Close close the server
func (p *Proxy) Close() error {
	return p.Server.Close()
}
