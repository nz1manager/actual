package main

import (
	"database/sql"
	"ielts-actual-backend/auth"
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

	// 1. Asosiy sahifa
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "IELTS Actual 2026 Backend Online"})
	})

	// 2. Google Login
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

	// 3. Foydalanuvchi ma'lumotlarini saqlash
	router.POST("/api/save-user", func(c *gin.Context) {
		var input struct {
			GoogleID  string `json:"google_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Phone     string `json:"phone"`
		}
		if err := c.ShouldBindJSON(&input); err != nil || input.FirstName == "" || input.LastName == "" {
			c.JSON(400, gin.H{"error": "Ma'lumotlar to'liq emas"})
			return
		}
		DB.Exec("UPDATE users SET first_name=$1, last_name=$2, phone=$3 WHERE google_id=$4",
			input.FirstName, input.LastName, input.Phone, input.GoogleID)
		c.JSON(200, gin.H{"message": "Updated successfully!"})
	})

	// 4. Admin Login
	router.POST("/api/admin/login", func(c *gin.Context) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		c.ShouldBindJSON(&creds)
		if creds.Username == "drkz1manager" && creds.Password == "4900728Tt$$$" {
			c.JSON(200, gin.H{"role": "admin", "token": "admin-secret"})
		} else {
			c.JSON(401, gin.H{"error": "Xato!"})
		}
	})

	// 5. O'quvchilarni boshqarish
	router.GET("/api/admin/all-students", func(c *gin.Context) {
		rows, _ := DB.Query("SELECT id, first_name, last_name, phone, email FROM users WHERE role = 'user'")
		defer rows.Close()
		var students []map[string]interface{}
		for rows.Next() {
			var id int
			var fn, ln, ph, em string
			rows.Scan(&id, &fn, &ln, &ph, &em)
			students = append(students, map[string]interface{}{"id": id, "name": fn + " " + ln, "phone": ph, "email": em})
		}
		c.JSON(200, students)
	})

	router.DELETE("/api/admin/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		DB.Exec("DELETE FROM users WHERE id = $1", id)
		c.JSON(200, gin.H{"message": "O'chirildi"})
	})

	// 6. Test natijalarini topshirish
	router.POST("/api/submit-test", func(c *gin.Context) {
		var sub struct {
			UserID int `json:"user_id"`; TestID int `json:"test_id"`; Correct int `json:"correct_answers"`; Total int `json:"total_questions"`
		}
		if err := c.ShouldBindJSON(&sub); err == nil {
			score := (float64(sub.Correct) / float64(sub.Total)) * 9
			DB.Exec("INSERT INTO scores (user_id, test_id, correct_answers, total_questions, score_value) VALUES ($1, $2, $3, $4, $5)",
				sub.UserID, sub.TestID, sub.Correct, sub.Total, score)
			c.JSON(200, gin.H{"message": "Test topshirildi", "score": score})
		}
	})

	// 7. Dashboard (Overview)
	router.GET("/api/user-dashboard/:id", func(c *gin.Context) {
		userID := c.Param("id")
		var total int; var avg float64
		DB.QueryRow("SELECT COUNT(*), COALESCE(AVG(score_value), 0) FROM scores WHERE user_id = $1", userID).Scan(&total, &avg)
		
		overview := "Qoniqarsiz"
		if avg >= 8.5 { overview = "A'lochi (Expert)" } else if avg >= 7.0 { overview = "Yaxshi (Good)" } else if avg >= 5.0 { overview = "Qoniqarli" }
		
		c.JSON(200, gin.H{"total_tests": total, "avg_score": avg, "overview": overview})
	})

	// --- FAQAT SHU BUYRUQ OXIRIDA TURISHI KERAK ---
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	router.Run(":" + port)
}
