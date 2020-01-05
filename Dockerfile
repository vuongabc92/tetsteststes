FROM golang:latest as builder

WORKDIR /app

COPY . .

COPY go.mod go.sum ./

RUN go mod download

RUN go build -o bin/octocv cmd/server/main.go