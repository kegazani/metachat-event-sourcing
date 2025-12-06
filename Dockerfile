FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -a -installsuffix cgo .

# This is a library, so we don't need a runtime stage
# The image is used for building other services
FROM scratch AS export
COPY --from=builder /app /