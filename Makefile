.bin/docker-proxy: $(wildcard *.go)
	@mkdir .bin
	GOOS=linux GOARCH=amd64 go build -o .bin/docker-proxy .

docker: .bin/docker-proxy
	docker build -t keegancsmith/docker-proxy .

.PHONY: docker test
