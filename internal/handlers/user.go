package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "yourmodule/internal/models"
    "yourmodule/internal/repository"
    "yourmodule/internal/utils"
)

type UserHandler struct {
    userRepo       *repository.UserRepository
    submissionRepo *repository.SubmissionRepository
}

func NewUserHandler(userRepo *repository.UserRepository, submissionRepo *repository.SubmissionRepository) *UserHandler {
    return &UserHandler{
        userRepo:       userRepo,
        submissionRepo: submissionRepo,
    }
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    userID := c.GetString("user_id")
    
    user, err := h.userRepo.FindByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    userID := c.GetString("user_id")
    
    var input struct {
        FirstName string `json:"first_name" binding:"required"`
        LastName  string `json:"last_name" binding:"required"`
        Phone     string `json:"phone"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "First name and last name are required"})
        return
    }
    
    user, err := h.userRepo.FindByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    user.FirstName = input.FirstName
    user.LastName = input.LastName
    user.Phone = input.Phone
    
    if err := h.userRepo.Update(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Updated successfully!",
        "user": user,
    })
}

func (h *UserHandler) GetStats(c *gin.Context) {
    userID := c.GetString("user_id")
    
    submissions, err := h.submissionRepo.FindByUserID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
        return
    }
    
    totalTests := len(submissions)
    totalScore := 0
    for _, sub := range submissions {
        if sub.Score != nil {
            totalScore += *sub.Score
        }
    }
    
    c.JSON(http.StatusOK, gin.H{
        "total_tests_taken": totalTests,
        "total_score":       totalScore,
        "average_score":     totalTests > 0 ? totalScore / totalTests : 0,
    })
}
