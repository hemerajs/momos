package momos

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ssiURL string

func SSITimeoutHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1000 * time.Millisecond)
	w.Header().Add("Cache-Control", "max-age=10")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>hello</h1>`)
}

func SSIInvalidStausCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Cache-Control", "max-age=10")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>hello</h1>`)
}

func SSIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=10")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>hello</h1>`)
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=10")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<ssi name="basket" timeout="200" template="true" src="`+ssiURL+`">Default content!	
		<ssi-timeout>
		<span>Please try it again!</span>
		</ssi-timeout>
		<ssi-error>
		<span>Please call the support!</span>
		</ssi-error>
	</ssi>
	`)
}

func TestSSIContent(t *testing.T) {

	ssiServer := httptest.NewServer(http.HandlerFunc(SSIHandler))
	ssiURL = ssiServer.URL
	defer ssiServer.Close()

	tsAPI := httptest.NewServer(http.HandlerFunc(APIHandler))
	defer tsAPI.Close()

	p := New(tsAPI.URL)

	server := httptest.NewServer(p.Handler)
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	assert.Equal(t, res.StatusCode, 200, "should return 200")
	assert.Equal(t, "<html><head></head><body><h1>hello</h1>\n\t</body></html>", bodyString)
}

func TestError(t *testing.T) {

	ssiServer := httptest.NewServer(http.HandlerFunc(SSIInvalidStausCodeHandler))
	ssiURL = ssiServer.URL
	defer ssiServer.Close()

	tsAPI := httptest.NewServer(http.HandlerFunc(APIHandler))
	defer tsAPI.Close()

	p := New(tsAPI.URL)

	server := httptest.NewServer(p.Handler)
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	assert.Equal(t, res.StatusCode, 200, "should return 200")
	assert.Equal(t, "<html><head></head><body>\n\t\t<span>Please call the support!</span>\n\t\t\n\t</body></html>", bodyString)
}
