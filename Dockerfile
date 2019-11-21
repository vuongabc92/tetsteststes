FROM golang:latest as builder
LABEL maintainer = "Vuong Bui <vuongabc92@gmail.com>"
WORKDIR /app
COPY . .
RUN go build -o bin/octocv cmd/server/server.go
EXPOSE 8080
CMD ["./octocv"]