package repository

import (
    "gorm.io/gorm"
    "yourmodule/internal/models"
)

type SubmissionRepository struct {
    db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) *SubmissionRepository {
    return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Create(submission *models.Submission) error {
    return r.db.Create(submission).Error
}

func (r *SubmissionRepository) FindByID(id string) (*models.Submission, error) {
    var submission models.Submission
    err := r.db.Preload("Test").Where("id = ?", id).First(&submission).Error
    return &submission, err
}

func (r *SubmissionRepository) FindByUserID(userID string) ([]models.Submission, error) {
    var submissions []models.Submission
    err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&submissions).Error
    return submissions, err
}

func (r *SubmissionRepository) FindByUserIDWithDetails(userID string) ([]models.Submission, error) {
    var submissions []models.Submission
    err := r.db.Preload("Test").Where("user_id = ?", userID).Order("created_at desc").Find(&submissions).Error
    return submissions, err
}

func (r *SubmissionRepository) FindByTestID(testID string) ([]models.Submission, error) {
    var submissions []models.Submission
    err := r.db.Where("test_id = ?", testID).Find(&submissions).Error
    return submissions, err
}
