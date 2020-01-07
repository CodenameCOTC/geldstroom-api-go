FROM golang:alpine AS build

WORKDIR /app/server
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build ./cmd/server/main.go

EXPOSE 4000

CMD ["./main"]
