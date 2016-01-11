package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var (
		sockPath string
		certPath string
	)
	flag.StringVar(&sockPath, "sockpath", "/var/run/docker.sock", "The path to the docker unix socket")
	flag.StringVar(&certPath, "certpath", "/certs", "Where to find or generate $DOCKER_CERT_PATH compatible certificates")
	flag.Parse()
	hosts := flag.Args()
	hosts = append(hosts, "127.0.0.1")

	serverCert, err := getOrGenerateServerCert(certPath, hosts)
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig, err := serverCert.TLSConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Generate client certs so the user can use them from the bindmount
	_, err = getOrGenerateClientCert(certPath)
	if err != nil {
		log.Fatal(err)
	}

	proxy := UnixSocketReverseProxy(sockPath)
	server := &http.Server{
		Addr:      ":2376",
		Handler:   proxy,
		TLSConfig: tlsConfig,
	}
	log.Fatal(server.ListenAndServeTLS("", ""))
}
