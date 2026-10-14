FROM golang:1.26.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /auth-service \
    ./cmd

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /auth-service /app/auth-service

EXPOSE 50051

ENTRYPOINT ["/app/auth-service"]
