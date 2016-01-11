package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
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
	proxy := unixSocketReverseProxy(sockPath)
	log.Fatal(http.ListenAndServe("127.0.0.1:10810", proxy))
}

func unixSocketReverseProxy(socketPath string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director:  func(_ *http.Request) {},
		Transport: &unixRoundTripper{socketPath},
	}
}

type unixRoundTripper struct {
	path string
}

func (u *unixRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	conn, err := net.Dial("unix", u.path)
	if err != nil {
		return nil, err
	}
	socketClientConn := httputil.NewClientConn(conn, nil)
	defer socketClientConn.Close()
	return socketClientConn.Do(req)
}
