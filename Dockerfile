# Build stage
FROM golang:1.26-alpine AS base
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download


FROM base as development
RUN go install github.com/air-verse/air@latest
CMD ["air", "-build.cmd", "go build -o ./tmp/main ./cmd/server", "-build.bin", "./tmp/main"]