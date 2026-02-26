FROM golang:1.21-alpine

WORKDIR /app

# Hammasini copy qil
COPY . .

# Build qil
RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
