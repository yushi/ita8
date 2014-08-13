package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"log"
)

var (
	addr string
	port = "4567"
)

func ita8br() {
	flag.Parse()
	addr = flag.Arg(0)
	if addr == "" {
		log.Fatal("address required")
	}

	dst := fmt.Sprintf("http://%s:%s", addr, port)
	dstURL, _ := url.Parse(dst)
	proxyHandler := httputil.NewSingleHostReverseProxy(dstURL)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: proxyHandler,
	}
	log.Fatal(server.ListenAndServe())
}

func ita8copy() {
	resp, err := http.Post("http://127.0.0.1:4567/", "text/plain", os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	println(string(b))
}

func ita8paste() {
	resp, err := http.Get("http://127.0.0.1:4567/")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	println(string(b))
}

func main() {
	if strings.HasSuffix(os.Args[0], "ita8br") {
		ita8br()
	} else if strings.HasSuffix(os.Args[0], "ita8copy") {
		ita8copy()
	} else if strings.HasSuffix(os.Args[0], "ita8paste") {
		ita8paste()
	}
}
