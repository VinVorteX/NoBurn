# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Build API and Worker binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker cmd/worker/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/worker .
COPY --from=builder /go/bin/migrate .
COPY migrations ./migrations
COPY entrypoint.sh .

# Make entrypoint executable
RUN chmod +x entrypoint.sh

# Expose API port
EXPOSE 3000

# Use entrypoint for auto-migrations
ENTRYPOINT ["./entrypoint.sh"]

# Default command (can be overridden in docker-compose)
CMD ["./api"]
