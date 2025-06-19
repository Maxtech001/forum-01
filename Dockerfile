# Base image
FROM alpine:3.17 AS base

# Install bash (useful for debugging or runtime scripts)
RUN apk add --no-cache bash

EXPOSE 8080

# === BUILD STAGE ===
FROM golang:1.19-alpine AS build

# Install necessary packages in one layer
RUN apk add --no-cache build-base sqlite

# Set working directory
WORKDIR /build

# Copy source code
COPY . .

# Create SQLite DB
RUN mkdir db && sqlite3 ./db/forum.db < db.sql

# Download Go dependencies
RUN go mod download

# Build Go app
RUN go build -o forum-docker .

# === FINAL STAGE ===
FROM base AS final

# Create app directory and set it as working dir
WORKDIR /app

# Copy built app and DB from build stage
COPY --from=build /build/forum-docker .
COPY --from=build /build/db ./db

# Run the app
CMD ["/app/forum-docker"]
