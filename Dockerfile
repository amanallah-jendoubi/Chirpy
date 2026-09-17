# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN go install github.com/air-verse/air@v1.52.3
COPY go.mod go.sum* ./
RUN go mod download
CMD ["air"]