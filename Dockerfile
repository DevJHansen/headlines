# Set the Go version to use
ARG GO_VERSION=1

# First stage: build the Go application
FROM golang:${GO_VERSION}-bookworm as builder

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -v -o /run-app .

# Second stage: create a minimal runtime environment
FROM debian:bookworm

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates

# Set the environment variable for Google Application Credentials
ENV GOOGLE_APPLICATION_CREDENTIALS=/usr/local/bin/firebase-sa.json

# Copy the built application and the service account key from the builder stage
COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/firebase-sa.json /usr/local/bin/firebase-sa.json

# Set the entry point for the container
CMD ["run-app"]
