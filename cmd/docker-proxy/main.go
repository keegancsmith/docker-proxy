package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"regexp"

	dproxy "github.com/keegancsmith/docker-proxy"
)

func main() {
	var (
		sockPath string
		certPath string
	)
	flag.StringVar(&sockPath, "sockpath", "/var/run/docker.sock", "The path to the docker unix socket")
	flag.StringVar(&certPath, "certpath", "/certs", "Where to find or generate $DOCKER_CERT_PATH compatible certificates")
	flag.Parse()
	suppliedHosts := flag.Args()

	hosts, err := ipsFromInterfaces()
	if err != nil {
		log.Fatal(err)
	}
	hosts = append(hosts, suppliedHosts...)

	serverCert, err := dproxy.GetOrGenerateServerCert(certPath, hosts)
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig, err := serverCert.TLSConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Generate client certs so the user can use them from the bindmount
	_, err = dproxy.GetOrGenerateClientCert(certPath)
	if err != nil {
		log.Fatal(err)
	}

	proxy := dproxy.UnixSocketReverseProxy(sockPath)
	server := &http.Server{
		Addr:      ":2376",
		Handler:   proxy,
		TLSConfig: tlsConfig,
	}
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func ipsFromInterfaces() ([]string, error) {
	ipRegexp := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
	ifts, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ips := []string{}
	for _, ift := range ifts {
		addrs, err := ift.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if addr.Network() == "ip+net" {
				if match := ipRegexp.FindString(addr.String()); match != "" {
					ips = append(ips, match)
				}
			}
		}
	}
	return ips, nil
}
