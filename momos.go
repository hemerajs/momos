package momos

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/hemerajs/momos/logger"
	"github.com/lox/httpcache"
)

const (
	version = "1.0.0"
	website = "https://github.com/hemerajs/momos"
	banner  = `
__  ___             
/  |/  /__  __ _  ___  ___
/ /|_/ / _ \/  ' \/ _ \(_-<
/_/  /_/\___/_/_/_/\___/___/ %s
Reverse proxy for advanced SSI									
`
)

var ServerLogging = false
var Log = logger.NewStdLogger(true, true, true, true, false)

type Proxy struct {
	ReverseProxy *httputil.ReverseProxy
	server       *http.Server
	ProxyURL     string
	Handler      *httpcache.Handler
}

func New(targetUrl string) *Proxy {
	target, tErr := url.Parse(targetUrl)

	if tErr != nil {
		Log.Errorf("Invalid url: %v", targetUrl)
		panic(tErr)
	}

	p := &Proxy{}
	p.ReverseProxy = httputil.NewSingleHostReverseProxy(target)
	p.ReverseProxy.Transport = &proxyTransport{http.DefaultTransport}

	return p
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p *Proxy) Start(addr string) error {
	// create proxy server
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	p.server = server
	// assign roundTrip handler
	p.server.Handler = p.ReverseProxy

	fmt.Printf(banner, version)
	fmt.Printf("server started on %s\n", server.Addr)

	// start server
	return server.ListenAndServe()
}
