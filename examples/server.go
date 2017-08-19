package main

import (
	"net/http"

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
	w.Header().Add("Cache-Control", "max-age=10")
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "examples/ssi.html")
}

func main() {
	// Api Mock
	api := http.NewServeMux()
	api.HandleFunc("/", backendHandler)
	go func() {
		err := http.ListenAndServe("127.0.0.1:8080", api)
		if err != nil {
			panic(err)
		}
	}()

	// SSI Mock
	ssi := http.NewServeMux()
	ssi.HandleFunc("/", ssiHandler)
	go func() {
		err := http.ListenAndServe("127.0.0.1:8081", ssi)
		if err != nil {
			panic(err)
		}
	}()

	momos.DebugLogging = true
	momos.ServerLogging = false

	// create momos instance and pass the url to your service
	p := momos.New("http://127.0.0.1:8080")
	// create proxy server
	server := &http.Server{Addr: "127.0.0.1:9090"}
	// assign roundTrip handler
	server.Handler = p.Handler
	// start server
	server.ListenAndServe()
}
