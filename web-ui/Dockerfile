# Build stage
FROM golang:1.21-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

# Set working directory to match repository structure
WORKDIR /repo

# Copy entire repository to maintain structure
COPY . .

# Build the web-ui application
WORKDIR /repo/web-ui
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-ui .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests (needed for AI services) and wget for health checks
RUN apk --no-cache add ca-certificates git wget

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory to match development structure
WORKDIR /repo

# Copy entire repository structure from builder
COPY --from=builder --chown=appuser:appgroup /repo .

# Make scripts executable
RUN find . -name "*.sh" -exec chmod +x {} \; 2>/dev/null || true

# Switch to non-root user
USER appuser

# Set working directory to web-ui for execution (maintains ../challenge-* paths)
WORKDIR /repo/web-ui

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./web-ui"]
