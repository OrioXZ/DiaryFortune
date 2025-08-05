# Stage 1: Build the Go binary
FROM golang:1.24.5 AS builder

WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# REMOVE this if present
# COPY fortunes.db .


# Enable CGO and build for Linux
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Install gcc for CGO (required for go-sqlite3)
RUN apt-get update && apt-get install -y gcc libc6-dev
RUN go build -o dairyfortune .

# Stage 2: Final image using Debian slim (not Alpine)
FROM debian:bookworm-slim

WORKDIR /root/

# Copy the built binary
COPY --from=builder /app/dairyfortune .

# Copy your database if needed
# COPY fortunes.db .    ← we can enable this later

EXPOSE 8080

CMD ["./dairyfortune"]
