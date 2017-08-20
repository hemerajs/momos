package momos

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/hemerajs/momos/logger"
	"github.com/lox/httpcache"
	"github.com/lox/httpcache/httplog"
)

var ServerLogging = false
var Log = logger.NewStdLogger(true, true, true, true, false)

type Proxy struct {
	ReverseProxy *httputil.ReverseProxy
	ProxyURL     string
	Handler      *httpcache.Handler
	Cache        httpcache.Cache
}

func PreCacheResponseHandler(h http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// https://github.com/lox/httpcache
		Log.Debugf("PreResponse (%v) - Cache is %v", req.Host+req.URL.String(), res.Header().Get("X-Cache"))
		h.ServeHTTP(res, req)
	}
}

func New(targetUrl string) *Proxy {
	target, tErr := url.Parse(targetUrl)

	if tErr != nil {
		Log.Errorf("Invalid url: %v", targetUrl)
		panic(tErr)
	}

	httpcache.DebugLogging = ServerLogging

	p := &Proxy{}
	p.ReverseProxy = httputil.NewSingleHostReverseProxy(target)
	p.ReverseProxy.Transport = &proxyTransport{http.DefaultTransport}

	p.Cache = httpcache.NewMemoryCache()
	p.Handler = httpcache.NewHandler(p.Cache, PreCacheResponseHandler(p.ReverseProxy))
	p.Handler.Shared = true

	respLogger := httplog.NewResponseLogger(p.Handler)
	respLogger.DumpRequests = true
	respLogger.DumpResponses = true
	respLogger.DumpErrors = true

	return p
}
