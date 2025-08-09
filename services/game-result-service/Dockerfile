# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy shared module first
COPY shared/ shared/

# Copy go mod files for this service
COPY services/game-result-service/go.mod services/game-result-service/go.sum services/game-result-service/

# Set working directory to service
WORKDIR /app/services/game-result-service

# Download dependencies (this will pull the shared module from local path)
RUN go mod download

# Copy service source code (shared already copied)
COPY services/game-result-service/ ./

# Build the application (we're already in the service directory)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/game-result-service/main .

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser

EXPOSE 8005

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8005/health || exit 1

CMD ["./main"]