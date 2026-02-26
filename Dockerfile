FROM golang:1.21-alpine

WORKDIR /app

# Faqat kerakli fayllarni nusxalash
COPY go.mod ./
# Agar go.sum bo'lsa uni ham qo'shing: COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
