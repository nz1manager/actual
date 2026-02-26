FROM golang:1.21-alpine

# Build uchun kerakli paketlarni o'rnatish
RUN apk add --no-cache git

WORKDIR /app

# Modul fayllarini nusxalash
COPY go.mod ./
# go.sum mavjud bo'lsa nusxalaydi
COPY go.sum* ./

# Kutubxonalarni yuklab olish
RUN go mod download

# Hamma kodni nusxalash
COPY . .

# Main.go'ni build qilish
RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
