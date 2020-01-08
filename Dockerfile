FROM golang:alpine AS build

WORKDIR /app/server
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

ENV GIN_MODE=release

RUN go build -ldflags="-s -w" ./cmd/server/main.go

EXPOSE 4000

CMD ["./main"]
