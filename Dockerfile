FROM golang:1.22.2-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY server ./server
COPY client ./client
COPY config ./config

FROM builder AS build

RUN go build -o dbuggen2 .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/dbuggen2 .

CMD ["/app/dbuggen2"]
