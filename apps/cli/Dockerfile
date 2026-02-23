# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the CLI
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cardwatcher-cli ./cmd/main.go

# Production stage
FROM alpine:3.20

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from build stage
COPY --from=build /app/cardwatcher-cli /app/cardwatcher-cli

# Run the CLI (entrypoint can be overridden)
ENTRYPOINT ["/app/cardwatcher-cli"]
