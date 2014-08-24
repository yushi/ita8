package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ita8paste() {
	_, body, err := req("GET", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(body)
}

func ita8copy() {
	_, body, err := req("PUT", os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(body)
}

func req(method string, r io.Reader) (resp *http.Response, body string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, "http://127.0.0.1:4567/", r)
	if err != nil {
		return
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	body = string(b)
	return
}
