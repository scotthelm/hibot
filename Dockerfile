FROM golang:latest

WORKDIR /go/src/github.com/scotthelm/hibot

RUN go get github.com/FiloSottile/gvt


COPY . /go/src/github.com/scotthelm/hibot

