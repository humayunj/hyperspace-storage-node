# syntax=docker/dockerfile:1

FROM golang:1.17.13-buster
# RUN apk add build-base


WORKDIR /app

COPY go.mod ./
COPY go.sum ./


COPY *.go ./
COPY config.yaml ./
COPY ./proto/* ./proto/
COPY ./contracts/* ./contracts/
COPY ./store/keystore ./store/keystore

RUN go mod download


RUN go build -o ./storage-node

# HTTP
EXPOSE 5555
# RPC
EXPOSE 8000 

ENV store_password=password

CMD ./storage-node
