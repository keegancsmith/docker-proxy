docker-proxy: $(wildcard *.go)
	GOOS=linux GOARCH=amd64 go build -o docker-proxy .

docker: docker-proxy
	docker build -t keegancsmith/docker-proxy .

.PHONY: docker
