FROM golang:1.25.1-alpine3.22 AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .