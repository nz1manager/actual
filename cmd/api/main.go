package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "yourmodule/internal/config"
    "yourmodule/internal/handlers"
    "yourmodule/internal/middleware"
    "yourmodule/internal/repository"
)

func main() {
    // Check if running in health check mode
    if len(os.Args) > 1 && os.Args[1] == "health" {
        // Simple health check - just exit with 0 if we can reach this point
        os.Exit(0)
    }

    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Initialize database
    db := config.InitDB()
    
    // Get database connection for health check
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get database connection:", err)
    }
    defer sqlDB.Close()

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

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        // Check database connection
        if err := sqlDB.Ping(); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "unhealthy",
                "database": "disconnected",
            })
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "database": "connected",
            "timestamp": time.Now().Unix(),
        })
    })

    // Configure CORS for Firebase app
    router.Use(middleware.CORS())

    // ... rest of your routes

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    router.Run(":" + port)
}
