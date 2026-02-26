package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "yourmodule/internal/config"
    "yourmodule/internal/handlers"
    "yourmodule/internal/middleware"
    "yourmodule/internal/repository"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Initialize database
    db := config.InitDB()
    
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    testRepo := repository.NewTestRepository(db)
    submissionRepo := repository.NewSubmissionRepository(db)

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(userRepo)
    userHandler := handlers.NewUserHandler(userRepo, submissionRepo)
    testHandler := handlers.NewTestHandler(testRepo, submissionRepo)
    adminHandler := handlers.NewAdminHandler(userRepo, testRepo, submissionRepo)

    // Setup router
    router := gin.Default()

    // Configure CORS for Firebase app
    router.Use(middleware.CORS())

    // Public routes
    auth := router.Group("/api/auth")
    {
        auth.GET("/google", authHandler.GoogleLogin)
        auth.GET("/google/callback", authHandler.GoogleCallback)
    }

    // Protected routes (require JWT)
    api := router.Group("/api")
    api.Use(middleware.JWTAuth())
    {
        // User routes
        api.GET("/user/profile", userHandler.GetProfile)
        api.PUT("/user/profile", userHandler.UpdateProfile)
        api.GET("/user/stats", userHandler.GetStats)
        
        // Test routes
        api.GET("/tests", testHandler.GetTests)
        api.GET("/tests/:id", testHandler.GetTest)
        api.POST("/tests/:id/submit", testHandler.SubmitTest)
        api.GET("/submissions/:id/review", testHandler.GetReview)
    }

    // Admin routes
    admin := router.Group("/api/admin")
    admin.Use(middleware.JWTAuth(), middleware.AdminOnly())
    {
        // Test management
        admin.POST("/tests", adminHandler.CreateTest)
        admin.PUT("/tests/:id", adminHandler.UpdateTest)
        admin.DELETE("/tests/:id", adminHandler.DeleteTest)
        admin.PATCH("/tests/:id/publish", adminHandler.TogglePublish)
        
        // Student management
        admin.GET("/students", adminHandler.GetStudents)
        admin.GET("/students/:id/results", adminHandler.GetStudentResults)
        admin.DELETE("/students/:id", adminHandler.DeleteStudent)
        
        // Leaderboard
        admin.GET("/leaderboard", adminHandler.GetLeaderboard)
    }

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    router.Run(":" + port)
}
