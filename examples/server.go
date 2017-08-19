package main

import (
	"html/template"
	"net/http"

	"github.com/hemerajs/momos"
)

type Page struct {
	Title string
}

func backendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=100000")
	w.Header().Set("Content-Type", "text/html")

	templ, err := template.ParseFiles("examples/advanced.html")

	if err != nil {
		panic(err)
	}

	templ.Execute(w, Page{Title: "Example SSI"})

}

func ssiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=100000")
	w.Header().Set("Content-Type", "text/html")

	templ, err := template.ParseFiles("examples/ssi.html")

	if err != nil {
		panic(err)
	}

	templ.Execute(w, Page{Title: "Example SSI"})

}

func main() {
	// API Mock
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

	p := momos.New("127.0.0.1:9090", "http://127.0.0.1:8080")
	err := p.Start()
	if err != nil {
		panic(err)
	}
}
