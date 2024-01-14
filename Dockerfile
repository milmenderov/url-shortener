FROM golang:1.21.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -a -installsuffix cgo -o urlshortener ./cmd/url-shortener/main.go
RUN mkdir -p /app/storage

FROM golang:1.21.5
WORKDIR /app

RUN mkdir -p /app/storage /app/config
COPY --from=builder /app/config /app/config
COPY --from=builder /app/urlshortener /app
CMD ["./urlshortener"]