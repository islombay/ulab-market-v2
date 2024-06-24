# Stage 1: Build stage
FROM golang:1.21.4 AS builder

WORKDIR /app


# Copy the entire source code
COPY . .

# Install make and other necessary packages
RUN apt-get update && apt-get install -y make

# Build the application
RUN make install

ENV ENV=docker
EXPOSE 8123

CMD ["make", "run"]