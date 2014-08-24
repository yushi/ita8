package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/yushi/ita8"
)

var (
	clipboardPath = "clipboard"
)

func ita8open(args []string) {
	b, err := json.Marshal(args)
	if err != nil {
		log.Fatal(err)
	}
	r := bytes.NewReader(b)
	if _, _, err := req(ita8.OpenPath, "POST", r); err != nil {
		log.Fatal(err)
	}
}

func ita8paste() {
	_, body, err := req(ita8.ClipboardPath, "GET", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(body)
}

func ita8copy() {
	_, body, err := req(ita8.ClipboardPath, "PUT", os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(body)
}

func req(path, method string, r io.Reader) (resp *http.Response, body string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("http://127.0.0.1:4567/%s", path), r)
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
