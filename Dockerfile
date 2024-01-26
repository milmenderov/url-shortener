FROM golang:1.21.5-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -a -installsuffix cgo -o urlshortener ./cmd/url-shortener/main.go


FROM golang:1.21.5-alpine3.19
WORKDIR /app

COPY --from=builder /app/config /app/config
COPY --from=builder /app/urlshortener /app
CMD ["./urlshortener"]