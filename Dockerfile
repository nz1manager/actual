FROM golang:1.21-alpine
WORKDIR /app
# Kutubxonalarni yuklash uchun mod fayllarni nusxalash
COPY go.mod go.sum ./
RUN go mod download
# Qolgan hamma narsani nusxalash
COPY . .
# Build qilish (xato bo'lsa shu yerda ko'rsatadi)
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
