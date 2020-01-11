FROM golang:alpine AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIN_MODE=release

WORKDIR /app/server    

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build ./cmd/server/main.go

WORKDIR /dist

RUN cp /app/server/main .

EXPOSE 4000

CMD ["/dist/main"]
