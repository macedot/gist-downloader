# Build stage
FROM golang:1.25-alpine AS builder

# Install git (required for cloning gists)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o gist-downloader .

# Runtime stage
FROM alpine:latest

# Install git in runtime stage as well
RUN apk add --no-cache git

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/gist-downloader .

# Create output directory
RUN mkdir -p /app/gist

# Set the binary as entrypoint
ENTRYPOINT ["./gist-downloader"]
