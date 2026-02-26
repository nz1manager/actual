FROM golang:1.21-alpine

WORKDIR /app

# Copy everything
COPY . .

# Disable checksum validation for this build
ENV GOINSECURE=gorm.io/*
ENV GONOSUMDB=gorm.io/*

# Build
RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
