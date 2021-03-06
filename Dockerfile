FROM golang:1.15 AS build

MAINTAINER imnothd "ideshenghe@gmail.com"

WORKDIR /go/src/github.com/IMNOTHD/go-kunpeng

COPY service service
COPY model model
COPY config config
COPY api api
COPY server server
COPY library library
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build .
RUN ls -a

FROM debian:latest AS prod

WORKDIR /app

COPY --from=build /go/src/github.com/IMNOTHD/go-kunpeng/go-kunpeng .
RUN ls -a

ENTRYPOINT ["/app/go-kunpeng"]