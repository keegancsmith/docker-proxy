package main

import (
	"flag"
	"log"

	dproxy "github.com/keegancsmith/docker-proxy"
)

func main() {
	var certPath string
	flag.StringVar(&certPath, "certpath", "/certs", "Where to find or generate $DOCKER_CERT_PATH compatible certificates")
	flag.Parse()
	_, err := dproxy.GetOrGenerateClientCert(certPath)
	if err != nil {
		log.Fatal(err)
	}
}
