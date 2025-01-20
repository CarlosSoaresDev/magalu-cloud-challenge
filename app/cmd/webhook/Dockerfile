# FROM golang:1.22-alpine as build
# ARG PORTS

# WORKDIR /app
# COPY app/go.mod app/go.sum ./
# RUN go mod download
# ADD ./app/. .
# COPY ./app/*.go ./
# RUN go build -o main main.go

# FROM alpine:latest as run
# WORKDIR /app
# COPY --from=build /app ./
# EXPOSE $PORT
# CMD ["/app/main"]


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
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /webhooks cmd/webhook/main.go

# Run stage
FROM alpine:3.19

COPY --from=build /webhooks /app/

# Expose API port
EXPOSE $PORT

CMD ["/app/webhooks"]