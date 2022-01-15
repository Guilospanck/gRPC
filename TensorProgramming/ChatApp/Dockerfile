FROM golang:alpine as build-env

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /chat_docker_grpc_example
RUN mkdir -p /chat_docker_grpc_example/proto

WORKDIR /chat_docker_grpc_example

COPY ./proto/service.pb.go /chat_docker_grpc_example/proto
COPY ./proto/service_grpc.pb.go /chat_docker_grpc_example/proto
COPY ./main.go /chat_docker_grpc_example

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build --ldflags "-s -w" -o chat_server .

# Multi-stage
FROM alpine
WORKDIR /project
RUN apk update && apk add --no-cache ca-certificates

COPY --from=build-env /chat_docker_grpc_example/chat_server /usr/bin/

ENTRYPOINT [ "/usr/bin/chat_server" ]