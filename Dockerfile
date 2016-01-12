FROM alpine:3.3

ADD .bin/* /usr/local/bin/

ENTRYPOINT ["docker-proxy"]

EXPOSE 2376
