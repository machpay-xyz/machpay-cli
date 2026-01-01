# ============================================================
# MachPay CLI - Docker Image
# ============================================================
#
# Simple image that copies pre-built binary from GoReleaser.
#
# Usage:
#   docker run --rm ghcr.io/machpay-xyz/cli version
#
# ============================================================

FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

# Copy pre-built binary from GoReleaser
COPY machpay /usr/local/bin/machpay

# Make executable
RUN chmod +x /usr/local/bin/machpay

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
