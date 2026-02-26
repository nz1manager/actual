package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/nz1manager/actual//internal/models"
    "github.com/nz1manager/actual//internal/repository"
)

type TestHandler struct {
    testRepo       *repository.TestRepository
    submissionRepo *repository.SubmissionRepository
}

func NewTestHandler(testRepo *repository.TestRepository, submissionRepo *repository.SubmissionRepository) *TestHandler {
    return &TestHandler{
        testRepo:       testRepo,
        submissionRepo: submissionRepo,
    }
}

func (h *TestHandler) GetTests(c *gin.Context) {
    testType := c.Query("type")
    
    var tests []models.Test
    var err error
    
    if testType == "" || testType == "All" {
        tests, err = h.testRepo.FindAllPublished()
    } else {
        tests, err = h.testRepo.FindByType(testType)
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tests"})
        return
    }
    
    c.JSON(http.StatusOK, tests)
}

func (h *TestHandler) GetTest(c *gin.Context) {
    testID := c.Param("id")
    
    test, err := h.testRepo.FindByID(testID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
        return
    }
    
    if !test.IsPublished {
        c.JSON(http.StatusForbidden, gin.H{"error": "Test not available"})
        return
    }
    
    // Remove answers from response
    response := struct {
        ID      string         `json:"id"`
        Title   string         `json:"title"`
        Type    string         `json:"type"`
        Content models.JSONMap `json:"content"`
    }{
        ID:      test.ID.String(),
        Title:   test.Title,
        Type:    test.Type,
        Content: test.Content,
    }
    
    c.JSON(http.StatusOK, response)
}

func (h *TestHandler) SubmitTest(c *gin.Context) {
    userID := c.GetString("user_id")
    testID := c.Param("id")
    
    var input struct {
        UserAnswers models.JSONMap `json:"user_answers" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission"})
        return
    }
    
    // Get test with answers
    test, err := h.testRepo.FindByIDWithAnswers(testID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
        return
    }
    
    // FIXED: Calculate score based on all test questions, not just answered ones
    score := calculateScore(input.UserAnswers, test.Answers, test.Content)
    
    // Create submission
    submission := &models.Submission{
        UserID:      uuid.MustParse(userID),
        TestID:      test.ID,
        UserAnswers: input.UserAnswers,
        Score:       &score,
    }
    
    if err := h.submissionRepo.Create(submission); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save submission"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Test submitted successfully",
        "score": score,
        "submission_id": submission.ID,
    })
}

func (h *TestHandler) GetReview(c *gin.Context) {
    submissionID := c.Param("id")
    userID := c.GetString("user_id")
    
    submission, err := h.submissionRepo.FindByID(submissionID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
        return
    }
    
    // Verify ownership
    if submission.UserID.String() != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
    
    // Get test with answers for comparison
    test, err := h.testRepo.FindByIDWithAnswers(submission.TestID.String())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load test"})
        return
    }
    
    // Create review data
    review := gin.H{
        "submission_id": submission.ID,
        "test_title":    test.Title,
        "test_type":     test.Type,
        "score":         submission.Score,
        "questions":     test.Content,
        "your_answers":  submission.UserAnswers,
        "correct_answers": test.Answers,
        "submitted_at":  submission.CreatedAt,
    }
    
    c.JSON(http.StatusOK, review)
}

// FIXED: Calculate score based on total questions from test content
func calculateScore(userAnswers, correctAnswers, testContent models.JSONMap) int {
    // Extract questions from test content structure
    // This assumes your test content has a "questions" array or map
    // Adjust this logic based on your actual test content structure
    
    var totalQuestions int
    
    // Method 1: If questions are stored as an array in the content
    if questions, ok := testContent["questions"].([]interface{}); ok {
        totalQuestions = len(questions)
    } else {
        // Method 2: If questions are stored as a map with question IDs as keys
        // Count the number of keys in correctAnswers (assuming each question has an answer)
        totalQuestions = len(correctAnswers)
    }
    
    if totalQuestions == 0 {
        return 0
    }
    
    // Count correct answers
    correct := 0
    for qID, userAns := range userAnswers {
        if correctAns, ok := correctAnswers[qID]; ok {
            // Handle different answer types (string, number, boolean, etc.)
            if compareAnswers(userAns, correctAns) {
                correct++
            }
        }
    }
    
    // Calculate percentage based on total questions
    return (correct * 100) / totalQuestions
}

// Helper function to compare answers of different types
func compareAnswers(userAnswer, correctAnswer interface{}) bool {
    switch v := correctAnswer.(type) {
    case string:
        if userStr, ok := userAnswer.(string); ok {
            return userStr == v
        }
    case float64:
        if userFloat, ok := userAnswer.(float64); ok {
            return userFloat == v
        }
    case bool:
        if userBool, ok := userAnswer.(bool); ok {
            return userBool == v
        }
    case []interface{}:
        // For multiple choice or array answers
        if userSlice, ok := userAnswer.([]interface{}); ok {
            if len(userSlice) != len(v) {
                return false
            }
            for i, val := range userSlice {
                if !compareAnswers(val, v[i]) {
                    return false
                }
            }
            return true
        }
    default:
        // Fallback to direct comparison
        return userAnswer == correctAnswer
    }
    return false
}
