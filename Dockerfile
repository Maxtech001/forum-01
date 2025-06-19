FROM alpine:3.17 as base

# To see inside the container workaround for audit
RUN apk add --no-cache bash

EXPOSE 8080

# Metainfo
#LABEL version="0.1" 
#LABEL status="dev"
#LABEL maintainers="Forum team"
#LABEL organization="kood/Jõhvi"

# BUILD
FROM golang:1.19-alpine AS build

# App settings
RUN mkdir /build

COPY . /build
WORKDIR /build

# Install C compiler and sqlite3
RUN apk add --no-cache build-base
RUN apk add --no-cache sqlite

# Database creation logic
RUN mkdir db && sqlite3 ./db/forum.db < db.sql

# Download dependencies
RUN go mod download

# Building go app
RUN go build -o forum-docker .

# DEPLOY
FROM base as final

# App settings
RUN mkdir /app
WORKDIR /app

# Copy build
COPY --from=build /build .

# Command to run the container
CMD ["/app/forum-docker"]
