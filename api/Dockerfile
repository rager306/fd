# Stage 1: Build static binary
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY . .

ARG VERSION=dev
ARG BUILD_HASH=unknown
ARG BUILD_DATE=unknown

# Build — static binary (no CGO, no libc deps)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildHash=${BUILD_HASH} -X main.BuildDate=${BUILD_DATE}" \
    -trimpath \
    -o /app/api .

# Stage 2: Runtime (minimal, static binary works on any base)
FROM debian:stable-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends curl \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /
COPY --from=builder /app/api /api
ENV GIN_MODE=release
ENTRYPOINT ["/api"]
