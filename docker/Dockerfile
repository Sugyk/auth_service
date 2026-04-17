ARG GO_VER=1.25.3
ARG ALPINE_VER=3.22

# Build stage
FROM golang:${GO_VER}-alpine${ALPINE_VER} AS builder

WORKDIR /src

# Install build dependencies
RUN apk add --no-cache \
    ca-certificates \
    git \
    curl

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with coverage support
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o app \
    -ldflags "-s -w -X main.version=${VERSION}" \
    ./cmd/auth_service

# Runtime stage
FROM alpine:${ALPINE_VER}

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Create non-root user
RUN addgroup -g 101 app && \
    adduser -H -u 101 -G app -s /bin/sh -D app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=app:app /src/app .

VOLUME ["/app/config.yaml", "/app/migrations"]

# Switch to non-root user
USER app

# Expose ports
EXPOSE 8080

# Run the application
CMD ["/app/app"]
