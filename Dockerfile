# Base stage with dependencies
FROM golang:1.24.1-alpine AS base
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Test stage
FROM base AS test
RUN go test ./... -v -cover

# Build stage
FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./src/cmd/api

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
