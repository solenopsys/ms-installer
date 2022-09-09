# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM  golang:1.19-buster

ARG TARGETARCH

WORKDIR /app/main/zmq_connector

COPY zmq_connector/*.go ./
COPY zmq_connector/go.mod ./
RUN go mod download

WORKDIR /app/main
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN go mod download

RUN  GOOS=linux GOARCH=$TARGETARCH go build  -o /go-binary

CMD [ "/go-binary" ]

