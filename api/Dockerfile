# Stage 1: Build static binary
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY . .

# Build — static binary (no CGO, no libc deps)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /app/api .

# Stage 2: Runtime (minimal, static binary works on any base)
FROM debian:stable-slim
WORKDIR /
COPY --from=builder /app/api /api
ENV GIN_MODE=release
ENTRYPOINT ["/api"]
