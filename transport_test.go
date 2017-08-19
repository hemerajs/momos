package momos

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "max-age=10")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<ssi
		name="basket"
		timeout="2000"
		template="true"
		src="http://localhost:8081">

		Default content!
		
		<ssi-timeout>
		<span>Please try it again! {{.DateLocal}}</span>
		</ssi-timeout>
		
		<ssi-error>
		<span>Please call the support!</span>
		</ssi-error>
	</ssi>
	`)
}

func TestMomos(t *testing.T) {

	tsApi := httptest.NewServer(http.HandlerFunc(apiHandler))
	defer tsApi.Close()

	p := New(tsApi.URL)

	server := httptest.NewServer(p.Handler)
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%v", string(body))

	assert.Equal(t, res.StatusCode, 200, "they should be equal")
}
