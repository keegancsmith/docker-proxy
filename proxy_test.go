package proxy

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUnixSocketReverseProxy(t *testing.T) {
	proxyTest(t, httptest.NewServer, http.Client{})
}

func TestUnixSocketTLSReverseProxy(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { os.RemoveAll(dir) }()

	serverCert, err := getOrGenerateServerCert(dir, []string{"127.0.0.1"})
	if err != nil {
		t.Fatal(err)
	}
	tlsConfig, err := serverCert.TLSConfig()
	if err != nil {
		t.Fatal(err)
	}

	startTestServer := func(h http.Handler) *httptest.Server {
		ts := httptest.NewUnstartedServer(h)
		ts.TLS = tlsConfig
		ts.StartTLS()
		return ts
	}

	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{Transport: httpTransport}

	proxyTest(t, startTestServer, client)
}

func proxyTest(t *testing.T, startTestServer func(http.Handler) *httptest.Server, client http.Client) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { os.RemoveAll(dir) }()

	sockPath := filepath.Join(dir, "docker.sock")
	cleanup, err := serveUnixHandler(sockPath)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// Then setup proxy and see if we successfully connect to the socket
	proxy := UnixSocketReverseProxy(sockPath)
	ts := startTestServer(proxy)
	defer ts.Close()

	res, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 404 {
		t.Fatal("Fetching top-level URL should 404")
	}

	res, err = client.Get(ts.URL + "/helloworld")
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if string(greeting) != "Hello World" {
		t.Fatal("Response body not Hello World", string(greeting))
	}
}

func serveUnixHandler(sockPath string) (func(), error) {
	// Setup an actual HTTP Handler on a unix socket
	l, err := net.Listen("unix", sockPath)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	go http.Serve(l, mux)
	return func() { l.Close() }, nil
}
