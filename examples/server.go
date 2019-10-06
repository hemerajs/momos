// +build ignore

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hemerajs/momos"
)

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
	apiAddr := "127.0.0.1:8080"
	ssiAddr := "127.0.0.1:8081"
	proxyAddr := "127.0.0.1:9090"

	// Api Mock
	api := http.NewServeMux()
	api.HandleFunc("/", backendHandler)
	go func() {
		http.ListenAndServe(apiAddr, api)
	}()

	// SSI Mock
	ssi := http.NewServeMux()
	ssi.HandleFunc("/", ssiHandler)
	go func() {
		http.ListenAndServe(ssiAddr, ssi)
	}()

	// Pass the url to your api
	p := momos.New("http://127.0.0.1:8080")

	// Start server
	go func() {
		if err := p.Start(proxyAddr); err != nil {
			fmt.Printf("shutting down the server, cause: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := p.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
