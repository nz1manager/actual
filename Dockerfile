FROM golang:1.21-alpine

WORKDIR /app

# Copy everything
COPY . .

# Disable checksum validation for ALL packages
ENV GOSUMDB=off
ENV GOINSECURE=*
ENV GOPRIVATE=*

# Build
RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
