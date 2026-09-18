# Build stage
FROM golang:1.22-alpine AS base
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download


FROM base as development
RUN go install github.com/air-verse/air@v1.52.3
CMD ["air"]
