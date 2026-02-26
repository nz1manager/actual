package main

import (
	"database/sql"
	"actual/auth" // <--- DIQQAT: go.mod dagi modul nomi bilan bir xil bo'lishi kerak
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func main() {
	godotenv.Load()

	var err error
	connStr := os.Getenv("DATABASE_URL")
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Bazaga ulanishda xato:", err)
	}

	router := gin.Default()

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

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "IELTS Actual 2026 Backend Online"})
	})

	router.POST("/api/google-login", func(c *gin.Context) {
		var body struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Token topilmadi"})
			return
		}
		googleUser, err := auth.GetGoogleUserInfo(body.Token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Google xatosi"})
			return
		}
		var userID int
		err = DB.QueryRow("SELECT id FROM users WHERE google_id = $1", googleUser.ID).Scan(&userID)
		if err != nil {
			err = DB.QueryRow(
				"INSERT INTO users (google_id, email, first_name, avatar_url) VALUES ($1, $2, $3, $4) RETURNING id",
				googleUser.ID, googleUser.Email, googleUser.Name, googleUser.Picture,
			).Scan(&userID)
		}
		c.JSON(200, gin.H{"user_id": userID, "full_name": googleUser.Name, "avatar": googleUser.Picture})
	})

	router.POST("/api/save-user", func(c *gin.Context) {
		var input struct {
			GoogleID  string `json:"google_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Phone     string `json:"phone"`
		}
		if err := c.ShouldBindJSON(&input); err == nil {
			DB.Exec("UPDATE users SET first_name=$1, last_name=$2, phone=$3 WHERE google_id=$4",
				input.FirstName, input.LastName, input.Phone, input.GoogleID)
			c.JSON(200, gin.H{"message": "Success"})
		}
	})

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	router.Run(":" + port)
}
