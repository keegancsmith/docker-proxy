.bin/docker-proxy: $(wildcard *.go) cmd/docker-proxy/main.go
	@mkdir -p .bin
	GOOS=linux GOARCH=amd64 go build -o .bin/docker-proxy ./cmd/docker-proxy

.bin/generate-client-certificate: $(wildcard *.go) cmd/generate-client-certificate/main.go
	@mkdir -p .bin
	GOOS=linux GOARCH=amd64 go build -o .bin/generate-client-certificate ./cmd/generate-client-certificate

docker: .bin/docker-proxy .bin/generate-client-certificate
	docker build -t keegancsmith/docker-proxy .

.PHONY: docker test
