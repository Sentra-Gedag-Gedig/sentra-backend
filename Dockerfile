FROM golang:1.21-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

# Install dependencies first (better caching)
RUN apk --no-cache add git

# Copy only what's needed for dependency resolution first
COPY go.mod go.sum ./
RUN go mod download

# Copy source files only (not logs, IDE configs, etc.)
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY database/ ./database/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Final stage with a smaller image and non-root user
FROM alpine:3.18

# Add necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -g '' -h /home/appuser appuser

WORKDIR /home/appuser

# Copy only what's needed for runtime
COPY --from=builder /app/main .
COPY --chown=appuser .env .

# Switch to the non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]