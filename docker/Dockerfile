FROM golang:1.11.3
RUN apt update
RUN apt-get install -y bash build-essential sqlite3 libsqlite3-dev vim
RUN go get github.com/toravir/glm
RUN mkdir -p /certs && cp /go/src/github.com/toravir/glm/*.pem /certs
RUN cp /go/src/github.com/toravir/glm/*.ini /tmp/
RUN cp /go/src/github.com/toravir/glm/tests/sanity.sh /go/bin/
EXPOSE 11223/tcp 
VOLUME ["/data", "/logs", "/certs"]
#ENTRYPOINT ["/go/bin/glm", "-config", "/tmp/example_cfg.ini"]
