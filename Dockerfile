# Build stage
FROM golang:1.21-alpine AS builder

# Install git and swag
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Run go mod tidy to ensure all dependencies are properly downloaded
RUN go mod tidy

RUN ls */**

# Generate swagger docs
RUN swag init

RUN go test ./... -v

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
