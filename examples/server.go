package main

import (
	"html/template"
	"net/http"

	"github.com/hemerajs/momos"
)

type Page struct {
	Title string
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=100000")
	w.Header().Set("Content-Type", "text/html")

	templ, err := template.ParseFiles("examples/test.html")

	if err != nil {
		panic(err)
	}

	templ.Execute(w, Page{Title: "Example SSI"})

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

	momos.DebugLogging = true
	// Start reverse proxy and replace "server" with "schmerver"
	p := momos.New("127.0.0.1:9090", "http://127.0.0.1:8080")
	err := p.Start()
	if err != nil {
		panic(err)
	}
}
