package main

import (
	"database/sql"
	"ielts-actual-backend/auth" // auth papkangiz borligini tekshiring
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func main() {
	// 1. .env yuklash
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

	// 4. CORS ruxsati
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

	// --- YO'NALISHLAR (ROUTES) BOSHLANDI ---

	// Test API
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "IELTS Actual 2026 Backend Online"})
	})

	// Google Login
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

	// User Save
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
		_, err = DB.Exec("UPDATE users SET first_name=$1, last_name=$2, phone=$3 WHERE google_id=$4",
			input.FirstName, input.LastName, input.Phone, input.GoogleID)
		c.JSON(200, gin.H{"message": "Updated successfully!"})
	})

	// Admin Login
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

	// O'quvchilar ro'yxati (13 & 15-qadam birlashmasi)
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

    // O'quvchini o'chirish
    router.DELETE("/api/admin/user/:id", func(c *gin.Context) {
        id := c.Param("id")
        DB.Exec("DELETE FROM users WHERE id = $1", id)
        c.JSON(200, gin.H{"message": "O'chirildi"})
    })

	// --- YO'NALISHLAR TUGADI ---

	// 5. Serverni ishga tushirish
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Test natijasini saqlash va hisoblash
router.POST("/api/submit-test", func(c *gin.Context) {
    var submission struct {
        UserID         int   `json:"user_id"`
        TestID         int   `json:"test_id"`
        CorrectAnswers int   `json:"correct_answers"`
        TotalQuestions int   `json:"total_questions"`
    }

    if err := c.ShouldBindJSON(&submission); err != nil {
        c.JSON(400, gin.H{"error": "Ma'lumotlar xato"})
        return
    }

    // Band score hisoblash (IELTS standartida taxminan)
    // score_value = (correct / total) * 9
    scoreValue := (float64(submission.CorrectAnswers) / float64(submission.TotalQuestions)) * 9

    // Bazaga natijani yozish
    _, err := DB.Exec(`
        INSERT INTO scores (user_id, test_id, correct_answers, total_questions, score_value) 
        VALUES ($1, $2, $3, $4, $5)`,
        submission.UserID, submission.TestID, submission.CorrectAnswers, submission.TotalQuestions, scoreValue)

    if err != nil {
        c.JSON(500, gin.H{"error": "Natijani saqlashda xato"})
        return
    }

    c.JSON(200, gin.H{
        "message": "Test muvaffaqiyatli topshirildi!",
        "your_score": scoreValue,
    })
})

// O'quvchi uchun Dashboard ma'lumotlari (Overview)
router.GET("/api/user-dashboard/:id", func(c *gin.Context) {
    userID := c.Param("id")

    var totalTests int
    var avgScore float64

    // Umumiy testlar soni va o'rtacha ballni olish
    err := DB.QueryRow(`
        SELECT COUNT(*), COALESCE(AVG(score_value), 0) 
        FROM scores WHERE user_id = $1`, userID).Scan(&totalTests, &avgScore)

    if err != nil {
        c.JSON(500, gin.H{"error": "Ma'lumotlarni yuklashda xato"})
        return
    }

    // Siz aytgan "Overview" darajasini aniqlash
    overview := "Qoniqarsiz"
    if avgScore >= 8.5 {
        overview = "A'lochi (Expert)"
    } else if avgScore >= 7.0 {
        overview = "Yaxshi (Good)"
    } else if avgScore >= 5.0 {
        overview = "Qoniqarli (Satisfactory)"
    }

    c.JSON(200, gin.H{
        "total_tests_done": totalTests,
        "total_score":      avgScore,
        "overview":         overview,
    })
})

	router.Run(":" + port)
}
