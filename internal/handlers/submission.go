package handlers

import (
    "github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
    // repo fields here
}

func NewSubmissionHandler() *SubmissionHandler {
    return &SubmissionHandler{}
}

func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
    c.JSON(200, gin.H{"message": "submission handler"})
}
