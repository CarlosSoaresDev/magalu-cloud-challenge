# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build

# Create work directory
WORKDIR /app

# Copy and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy API packages
COPY . ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /api cmd/api/main.go

# Run stage
FROM alpine:3.19

WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build /api /app/

# Expose API port
EXPOSE $PORT

CMD ["/app/api"]
