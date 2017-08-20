package momos

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/hemerajs/momos/logger"
	"github.com/lox/httpcache"
)

var ServerLogging = false
var Log = logger.NewStdLogger(true, true, true, true, false)

type Proxy struct {
	ReverseProxy *httputil.ReverseProxy
	ProxyURL     string
	Handler      *httpcache.Handler
}

// PreCacheResponseHandler is an http handler to log informations about the cache status
// https://github.com/lox/httpcache
func PreCacheResponseHandler(h http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		url := req.Host + req.URL.String()
		cacheHeader := res.Header().Get("X-Cache")
		Log.Noticef("PreResponse url: %v, cache: %v", url, cacheHeader)
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

	cache := httpcache.NewMemoryCache()
	p.Handler = httpcache.NewHandler(cache, PreCacheResponseHandler(p.ReverseProxy))
	p.Handler.Shared = true

	return p
}
