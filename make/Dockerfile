FROM golang:1.14.6 AS builder

RUN mkdir -p /go/src/github.com/congcongke/httpfileserver

WORKDIR /go/src/github.com/congcongke/httpfileserver

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLE=0

RUN go build -o bin/httpserver cmd/main.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/congcongke/httpfileserver/bin/httpserver .

ENTRYPOINT ["./httpserver"]
