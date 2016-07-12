FROM golang:1.6.2

RUN apt-get update && apt-get -y install git mysql-client

RUN mkdir -p /go && mkdir -p /root/app

ENV GOPATH /go
ENV HOME /root
WORKDIR ${HOME}

RUN go get github.com/mattn/gom

WORKDIR /root/app
CMD ["/bin/bash"]
