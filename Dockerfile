# Stage 1: Build the Go app
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
# Download dependencies
RUN go mod tidy
# Build the binary
RUN go build -o sentinel

# Stage 2: Create a tiny image to run it
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sentinel .
CMD ["./sentinel"]