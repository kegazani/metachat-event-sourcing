FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the library
RUN go build -a -installsuffix cgo .

# This is a library, so we don't need a runtime stage
# The image is used for building other services
FROM scratch AS export
COPY --from=builder /app /