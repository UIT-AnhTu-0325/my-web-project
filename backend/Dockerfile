# Start from the official Golang image as a build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app (adjust the path to your main.go if needed)
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Start a minimal image for running
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/server ./server

# Expose the port your app runs on
EXPOSE 8080

# Run the server binary
ENTRYPOINT ["/app/server"]
