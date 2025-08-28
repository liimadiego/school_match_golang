FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod file
COPY go.mod ./

# Download dependencies
RUN go mod download && go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o schoolmatch .

# Create final lightweight image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder image
COPY --from=builder /app/schoolmatch .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./schoolmatch"]