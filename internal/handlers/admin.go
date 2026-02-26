package handlers

import (
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/nz1manager/actual/internal/models"
    "github.com/nz1manager/actual/internal/repository"
)

type AdminHandler struct {
    userRepo       *repository.UserRepository
    testRepo       *repository.TestRepository
    submissionRepo *repository.SubmissionRepository
}

func NewAdminHandler(
    userRepo *repository.UserRepository,
    testRepo *repository.TestRepository,
    submissionRepo *repository.SubmissionRepository,
) *AdminHandler {
    return &AdminHandler{
        userRepo:       userRepo,
        testRepo:       testRepo,
        submissionRepo: submissionRepo,
    }
}

// Admin login handler (separate from Google OAuth)
func (h *AdminHandler) AdminLogin(c *gin.Context) {
    var credentials struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    
    // Check credentials against environment variables
    if credentials.Username != os.Getenv("ADMIN_USERNAME") || 
       credentials.Password != os.Getenv("ADMIN_PASSWORD") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Generate admin token
    token, err := utils.GenerateToken("admin", true)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"token": token})
}

// Test Management
func (h *AdminHandler) CreateTest(c *gin.Context) {
    var test models.Test
    if err := c.ShouldBindJSON(&test); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test data"})
        return
    }
    
    if err := h.testRepo.Create(&test); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test"})
        return
    }
    
    c.JSON(http.StatusCreated, test)
}

func (h *AdminHandler) UpdateTest(c *gin.Context) {
    testID := c.Param("id")
    
    var updates models.Test
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
        return
    }
    
    test, err := h.testRepo.FindByID(testID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
        return
    }
    
    // Update fields
    test.Title = updates.Title
    test.Type = updates.Type
    test.Content = updates.Content
    test.Answers = updates.Answers
    
    if err := h.testRepo.Update(test); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Test updated successfully", "test": test})
}

func (h *AdminHandler) DeleteTest(c *gin.Context) {
    testID := c.Param("id")
    
    if err := h.testRepo.Delete(testID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Test deleted successfully"})
}

func (h *AdminHandler) TogglePublish(c *gin.Context) {
    testID := c.Param("id")
    
    test, err := h.testRepo.FindByID(testID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
        return
    }
    
    test.IsPublished = !test.IsPublished
    
    if err := h.testRepo.Update(test); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Test publish status updated",
        "is_published": test.IsPublished,
    })
}

// Student Management
func (h *AdminHandler) GetStudents(c *gin.Context) {
    students, err := h.userRepo.FindAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
        return
    }
    
    // Get scores for each student
    var result []gin.H
    for _, student := range students {
        submissions, _ := h.submissionRepo.FindByUserID(student.ID.String())
        totalScore := 0
        for _, sub := range submissions {
            if sub.Score != nil {
                totalScore += *sub.Score
            }
        }
        
        result = append(result, gin.H{
            "id":         student.ID,
            "name":       student.FirstName + " " + student.LastName,
            "email":      student.Email,
            "phone":      student.Phone,
            "total_score": totalScore,
            "tests_taken": len(submissions),
            "created_at": student.CreatedAt,
        })
    }
    
    c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) GetStudentResults(c *gin.Context) {
    studentID := c.Param("id")
    
    submissions, err := h.submissionRepo.FindByUserIDWithDetails(studentID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results"})
        return
    }
    
    c.JSON(http.StatusOK, submissions)
}

func (h *AdminHandler) DeleteStudent(c *gin.Context) {
    studentID := c.Param("id")
    
    // This will cascade delete all submissions due to foreign key constraint
    if err := h.userRepo.Delete(studentID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Student and all associated data deleted successfully"})
}

// Leaderboard
func (h *AdminHandler) GetLeaderboard(c *gin.Context) {
    users, err := h.userRepo.FindAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
        return
    }
    
    var leaderboard []gin.H
    for _, user := range users {
        submissions, _ := h.submissionRepo.FindByUserID(user.ID.String())
        
        totalScore := 0
        for _, sub := range submissions {
            if sub.Score != nil {
                totalScore += *sub.Score
            }
        }
        
        avgScore := 0
        if len(submissions) > 0 {
            avgScore = totalScore / len(submissions)
        }
        
        // Determine rank based on average score
        rank := "Beginner"
        if avgScore >= 80 {
            rank = "Advanced"
        } else if avgScore >= 60 {
            rank = "Intermediate"
        } else if avgScore >= 40 {
            rank = "Developing"
        }
        
        leaderboard = append(leaderboard, gin.H{
            "student_name": user.FirstName + " " + user.LastName,
            "total_score":  totalScore,
            "average_score": avgScore,
            "tests_taken":  len(submissions),
            "rank":         rank,
        })
    }
    
    c.JSON(http.StatusOK, leaderboard)
}
