package main

import (
	"net"
	"net/http"
	"net/http/httputil"
)

func UnixSocketReverseProxy(socketPath string) *httputil.ReverseProxy {
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
