= docker-daemon-proxy

Proxy to expose the docker unix socket over tcp. Run this with docker

    $ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock keegancsmith/docker-daemon-proxy 
