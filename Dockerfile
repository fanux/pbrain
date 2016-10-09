FROM golang:latest

COPY . /go/src/github.com/fanux/pbrain/

RUN go get github.com/tools/godep && cd /go/src/github.com/fanux/pbrain/ && godep go install
