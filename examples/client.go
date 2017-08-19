// +build ignore
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	client := http.Client{}
	resp, err := client.Get("http://127.0.0.1:9090/llamas")
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("%v\n", string(b))
	fmt.Printf("Cached: %v\n", resp.Header.Get("X-Cache"))

}
