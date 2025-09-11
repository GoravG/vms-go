# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Final stage
FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server /app/server

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["/app/server"]