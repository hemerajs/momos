package momos

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkCachingFiles(b *testing.B) {

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := http.Client{}
	defer server.Close()

	mux.HandleFunc("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<html><head></head><body>
			<ssi name="basket" timeout="200" template="true" src="`+server.URL+`/ssi">Default content!	
			<ssi-timeout>
			<span>Please try it again!</span>
			</ssi-timeout>
			<ssi-error>
			<span>Please call the support!</span>
			</ssi-error>
		</ssi>
		</body></html>
		`)
	}))

	mux.HandleFunc("/ssi", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<h1>hello</h1>`)
	}))

	p := New(server.URL)
	proxy := httptest.NewServer(p.ReverseProxy)
	defer proxy.Close()

	for n := 0; n < b.N; n++ {
		resp, err := client.Get(fmt.Sprintf("%s/api/%d", proxy.URL, n))
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}
