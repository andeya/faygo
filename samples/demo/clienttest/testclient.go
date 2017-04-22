package main

import (
	"io/ioutil"
	"net/http"
)

func main() {
	for i := 0; i < 500; i++ {
		go do()
	}
	select {}
}
func do() {
	for {
		resp, _ := http.Get("http://localhost:8080")
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
}
