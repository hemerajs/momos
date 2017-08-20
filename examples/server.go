// +build ignore

package main

import (
	"net/http"
	"time"

	"github.com/hemerajs/momos"
)

type Page struct {
	Title string
}

func backendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=10")
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "examples/advanced.html")
}

func ssiHandler(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(1000 * time.Millisecond)
	w.Header().Add("Cache-Control", "max-age=10")
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "examples/ssi.html")
}

func main() {
	// Api Mock
	api := http.NewServeMux()
	api.HandleFunc("/", backendHandler)
	go func() {
		http.ListenAndServe("127.0.0.1:8080", api)
	}()

	// SSI Mock
	ssi := http.NewServeMux()
	ssi.HandleFunc("/", ssiHandler)
	go func() {
		http.ListenAndServe("127.0.0.1:8081", ssi)
	}()

	// create momos instance and pass the url to your service
	p := momos.New("http://127.0.0.1:8080")
	// create proxy server
	server := &http.Server{
		Addr:         "127.0.0.1:9090",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// assign roundTrip handler
	server.Handler = p.Handler
	// start server
	server.ListenAndServe()
}
