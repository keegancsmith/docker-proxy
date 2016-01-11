FROM alpine:3.3

ADD .bin/docker-proxy /usr/local/bin/docker-proxy

ENTRYPOINT ["docker-proxy"]

EXPOSE 2376
