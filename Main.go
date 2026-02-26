package main

import (
	"log"
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"database/sql"
)

var DB *sql.DB

func main() {
	// 1. .env faylni yuklash
	godotenv.Load()

	// 2. Bazaga ulanish
	var err error
	connStr := os.Getenv("DATABASE_URL")
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Bazaga ulanishda xato:", err)
	}

	// 3. Serverni yaratish
	router := gin.Default()

	// 4. Dizayn shaffof chiqishi uchun Front-endga ruxsat (CORS)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Asosiy sahifa testi
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "IELTS Actual 2026 Backend Online"})
	})
// Google orqali kirish va bazaga yozish
router.POST("/api/google-login", func(c *gin.Context) {
    var body struct {
        Token string `json:"token"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(400, gin.H{"error": "Token topilmadi"})
        return
    }

    // Google'dan ma'lumotlarni olish (auth paketi orqali)
    googleUser, err := auth.GetGoogleUserInfo(body.Token)
    if err != nil {
        c.JSON(401, gin.H{"error": "Google bilan bog'lanib bo'lmadi"})
        return
    }

    // Bazada foydalanuvchini qidirish yoki yangi ochish
    var userID int
    err = DB.QueryRow("SELECT id FROM users WHERE google_id = $1", googleUser.ID).Scan(&userID)

    if err != nil { // Foydalanuvchi topilmadi, demak yangi qo'shamiz
        err = DB.QueryRow(
            "INSERT INTO users (google_id, email, first_name, avatar_url) VALUES ($1, $2, $3, $4) RETURNING id",
            googleUser.ID, googleUser.Email, googleUser.Name, googleUser.Picture,
        ).Scan(&userID)
    }

    c.JSON(200, gin.H{
        "user_id":    userID,
        "google_id":  googleUser.ID,
        "full_name":  googleUser.Name,
        "email":      googleUser.Email,
        "avatar":     googleUser.Picture,
        "is_new":     err == nil, // Agar yangi bo'lsa, front-end Ism/Tel kiritishga yo'naltiradi
    })
})

	// 5. Siz aytgan "Save" funksiyasi (Foydalanuvchi ma'lumotlarini saqlash)
	router.POST("/api/save-user", func(c *gin.Context) {
		var input struct {
			GoogleID  string `json:"google_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Phone     string `json:"phone"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Ma'lumotlar xato"})
			return
		}

		// Ism va Familiya bo'sh bo'lsa save tugmasi ishlamasligi uchun tekshiruv
		if input.FirstName == "" || input.LastName == "" {
			c.JSON(400, gin.H{"error": "Ism va Familiya majburiy!"})
			return
		}

		// Bazada yangilash (Update)
		_, err := DB.Exec("UPDATE users SET first_name=$1, last_name=$2, phone=$3 WHERE google_id=$4",
			input.FirstName, input.LastName, input.Phone, input.GoogleID)

		if err != nil {
			c.JSON(500, gin.H{"error": "Bazaga saqlashda xato"})
			return
		}

		c.JSON(200, gin.H{"message": "Updated successfully!"})
	})

	// Serverni yoqish
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	router.Run(":" + port)
}
