# ============================================================
# MachPay CLI - Docker Image
# ============================================================
#
# Multi-stage build for minimal image size.
#
# Usage:
#   docker build -t machpay/cli .
#   docker run --rm machpay/cli version
#
# ============================================================

# Build stage
FROM golang:1.22-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build arguments
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" \
    -trimpath \
    -o machpay \
    ./cmd/machpay

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

# Copy binary from builder
COPY --from=builder /build/machpay /usr/local/bin/machpay

# Create non-root user
RUN addgroup -g 1000 machpay && \
    adduser -u 1000 -G machpay -h /home/machpay -D machpay

# Create config directory
RUN mkdir -p /home/machpay/.machpay && \
    chown -R machpay:machpay /home/machpay

# Switch to non-root user
USER machpay
WORKDIR /home/machpay

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD machpay version || exit 1

# Default command
ENTRYPOINT ["machpay"]
CMD ["--help"]

