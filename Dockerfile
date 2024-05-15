FROM golang:1.22.2
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY server ./server
COPY client ./client
COPY assets ./assets

EXPOSE 8080
CMD ["go", "run", "main.go"]