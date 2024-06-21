# Stage 1: Build stage
FROM golang:1.21.4 AS builder

WORKDIR /app


# Copy the entire source code
COPY . .

# Install make and other necessary packages
RUN apt-get update && apt-get install -y make

# Build the application
RUN make install

# Stage 2: Runtime stage
FROM alpine:latest AS runtime

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app /app/bin/myapp

# Set environment variables
ENV ENV=docker
EXPOSE 8123

# Run the application
CMD ["/app/bin/myapp"]