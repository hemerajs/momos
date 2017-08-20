package momos

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var s struct {
	server *httptest.Server
	proxy  *httptest.Server
	client http.Client
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {

	mux := http.NewServeMux()
	s.server = httptest.NewServer(mux)
	s.client = http.Client{}

	mux.HandleFunc("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<html><head></head><body>
			<ssi name="basket" timeout="200" template="true" src="`+s.server.URL+`/ssi">Default content!	
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

	mux.HandleFunc("/api/ssi500", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<ssi name="basket" timeout="200" template="true" src="`+s.server.URL+`/ssi500">Default content!	
			<ssi-timeout>
			<span>Please try it again!</span>
			</ssi-timeout>
			<ssi-error>
			<span>Please call the support!</span>
			</ssi-error>
		</ssi>
		`)
	}))

	mux.HandleFunc("/api/ssiTimeout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<ssi name="basket" timeout="500" template="true" src="`+s.server.URL+`/ssiTimeout">Default content!	
			<ssi-timeout>
			<span>Please try it again!</span>
			</ssi-timeout>
			<ssi-error>
			<span>Please call the support!</span>
			</ssi-error>
		</ssi>
		`)
	}))

	mux.HandleFunc("/api/filterIncludes", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
			<ssi name="basket" timeout="500" template="true" src="`+s.server.URL+`/withIncludes">Default content!	
			<ssi-timeout>
			<span>Please try it again!</span>
			</ssi-timeout>
			<ssi-error>
			<span>Please call the support!</span>
			</ssi-error>
		</ssi>
		`)
	}))

	p := New(s.server.URL)
	s.proxy = httptest.NewServer(p.Handler)

	mux.HandleFunc("/ssi", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<h1>hello</h1>`)
	}))

	mux.HandleFunc("/ssiTimeout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(600 * time.Millisecond)
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<h1>hello</h1>`)
	}))

	mux.HandleFunc("/ssi500", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<h1>hello</h1>`)
	}))

	mux.HandleFunc("/withIncludes", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=10")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<h1>hello</h1><script></script><link></link>`)
	}))
}

func teardown() {
	s.server.Close()
	s.proxy.Close()
}

var ssiURL string

func TestSSIContent(t *testing.T) {
	res, err := s.client.Get(s.proxy.URL + "/api")
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	assert.Equal(t, res.StatusCode, 200, "should return 200")
	assert.Equal(t, "<html><head></head><body>\n\t\t\t<h1>hello</h1>\n\t\t\n\t\t</body></html>", bodyString)
}

func TestError(t *testing.T) {
	res, err := s.client.Get(s.proxy.URL + "/api/ssi500")
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	assert.Equal(t, res.StatusCode, 200, "should return 200")
	assert.Equal(t, "<html><head></head><body>\n\t\t\t<span>Please call the support!</span>\n\t\t\t\n\t\t</body></html>", bodyString)
}

func TestTimeout(t *testing.T) {
	res, err := s.client.Get(s.proxy.URL + "/api/ssiTimeout")
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	assert.Equal(t, res.StatusCode, 200, "should return 200")
	assert.Equal(t, "<html><head></head><body>\n\t\t\t<span>Please try it again!</span>\n\t\t\t\n\t\t</body></html>", bodyString)
}

func TestFilterIncludes(t *testing.T) {
	res, err := s.client.Get(s.proxy.URL + "/api/filterIncludes")
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	assert.Equal(t, res.StatusCode, 200, "should return 200")
	assert.Equal(t, "<html><head></head><body><h1>hello</h1>\n\t\t</body></html>", bodyString)
}
