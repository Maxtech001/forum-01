# === BASE STAGE ===
FROM alpine:3.17 AS base

# Bash for debugging or interactive shell use
RUN apk add --no-cache bash

EXPOSE 8080

# === BUILD STAGE ===
FROM golang:1.19-alpine AS build

# Install build dependencies
RUN apk add --no-cache build-base sqlite

# Set up build directory
WORKDIR /build

# Copy everything into the container
COPY . .

# Create database
RUN mkdir db && sqlite3 ./db/forum.db < db.sql

# Download Go dependencies
RUN go mod download

# Build the Go application
RUN go build -o forum-docker .

# === FINAL STAGE ===
FROM base AS final

WORKDIR /app

# Copy the built binary
COPY --from=build /build/forum-docker .

# ✅ Copy required runtime folders
COPY --from=build /build/templates ./templates
COPY --from=build /build/styles ./styles
COPY --from=build /build/js ./js
COPY --from=build /build/db ./db

# Run the application
CMD ["./forum-docker"]

