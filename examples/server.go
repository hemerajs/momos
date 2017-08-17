package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/hemerajs/momos"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=100000")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<ssi-teaser>hello, you've hit the server %s\n %v </ssi-teaser>", r.URL.Path, rand.Int())
}

func main() {
	// API Mock
	api := http.NewServeMux()
	api.HandleFunc("/", hello)
	go func() {
		err := http.ListenAndServe("127.0.0.1:8080", api)
		if err != nil {
			panic(err)
		}
	}()

	// Start reverse proxy and replace "server" with "schmerver"
	p := momos.New("127.0.0.1:9090", "http://127.0.0.1:8080")
	err := p.Start()
	if err != nil {
		panic(err)
	}
}
