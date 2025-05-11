# syntax=docker/dockerfile:1

# Start from Go base image
FROM golang:1.24.3

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN go build -o /simulator

# Run
CMD ["/simulator"]
