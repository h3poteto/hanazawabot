FROM golang:1.7.3-alpine

RUN apk add --no-cache \
    git \
    build-base \
    openssh

RUN mkdir -p /go && \
    mkdir -p /root/app

ENV GOPATH /go
ENV HOME /root
WORKDIR ${HOME}

RUN go get github.com/mattn/gom
RUN go get -u github.com/jteeuwen/go-bindata/...

WORKDIR /root/app

CMD ["/bin/sh"]
