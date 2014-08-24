package main

import (
	"flag"
	"fmt"
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

func main() {
	if strings.HasSuffix(os.Args[0], "ita8br") {
		ita8br()
	} else if strings.HasSuffix(os.Args[0], "copy") {
		ita8copy()
	} else if strings.HasSuffix(os.Args[0], "paste") {
		ita8paste()
	}
}
