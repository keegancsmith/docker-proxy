package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/keegancsmith/docker-daemon-proxy"
)

func main() {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { os.RemoveAll(dir) }()

	// Setup HTTP Handler on a unix socket
	sockPath := filepath.Join(dir, "docker.sock")
	l, err := net.Listen("unix", sockPath)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	go http.Serve(l, mux)

	// Then setup proxy
	proxy := proxy.UnixSocketReverseProxy(sockPath)
	log.Fatal(http.ListenAndServe("127.0.0.1:10810", proxy))
}
