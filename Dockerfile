FROM golang:1.14 AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPATH=/go
WORKDIR /go/src/github.com/allanhung/ali-ecs-tag-update

COPY cmd cmd/
COPY pkg pkg/
COPY main.go go.mod ./

RUN GOARCH=amd64 GOOS=linux go build -o /go/bin/ali-ecs-tag-update main.go

FROM debian:stretch-slim
RUN apt-get update && apt-get install -y git

ENV XDG_CONFIG_HOME=/opt

COPY --from=builder /go/bin/ali-ecs-tag-update /usr/bin/ali-ecs-tag-update

WORKDIR /working
ENTRYPOINT ["/usr/bin/ali-ecs-tag-update", "updatek8stags"]
