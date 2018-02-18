FROM h3poteto/golang:1.9.4 AS package

ENV APPROOT /go/src/github.com/h3poteto/hanazawabot

COPY --chown=go:go . ${APPROOT}
WORKDIR ${APPROOT}

RUN set -x \
    && dep ensure \
    && go generate \
    && GOOS=linux GOARCH=amd64 go build -o pkg/hanazawabot -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"'


FROM alpine:3.7

ENV APPROOT /var/opt/app

COPY --from=package /go/src/github.com/h3poteto/hanazawabot/pkg ${APPROOT}
WORKDIR ${APPROOT}
