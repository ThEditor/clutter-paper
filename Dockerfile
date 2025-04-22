FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o clutter-paper ./cmd/app.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/clutter-paper .

EXPOSE 8080

CMD ["./clutter-paper"]
