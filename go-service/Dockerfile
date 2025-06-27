# Development Dockerfile
FROM golang:1.21

# Set working directory
WORKDIR /app

# Install air for hot reloading (use compatible version)
RUN go install github.com/cosmtrek/air@v1.49.0

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (will be overridden by volume in dev)
COPY . .

# Create reports directory
RUN mkdir -p reports

# Expose port
EXPOSE 8080

# Use air for hot reloading in development
CMD ["air", "-c", ".air.toml"]

