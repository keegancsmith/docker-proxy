# docker-proxy

Proxy to expose the docker unix socket over tcp. Run this with docker

    $ docker run --rm -p 10810:2376 \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v $PWD/certs:/certs \
    keegancsmith/docker-proxy 192.168.99.100
