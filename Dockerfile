# Super simple Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Faqat go.mod ni copy qilamiz
COPY go.mod ./
RUN go mod download

# Hammasini copy qilamiz
COPY . .

# Build qilamiz
RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
