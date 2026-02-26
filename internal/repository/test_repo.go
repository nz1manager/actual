package repository

import (
    "gorm.io/gorm"
    "yourmodule/internal/models"
)

type TestRepository struct {
    db *gorm.DB
}

func NewTestRepository(db *gorm.DB) *TestRepository {
    return &TestRepository{db: db}
}

func (r *TestRepository) Create(test *models.Test) error {
    return r.db.Create(test).Error
}

func (r *TestRepository) FindByID(id string) (*models.Test, error) {
    var test models.Test
    err := r.db.Where("id = ?", id).First(&test).Error
    return &test, err
}

func (r *TestRepository) FindByIDWithAnswers(id string) (*models.Test, error) {
    var test models.Test
    err := r.db.Select("*").Where("id = ?", id).First(&test).Error
    return &test, err
}

func (r *TestRepository) FindAllPublished() ([]models.Test, error) {
    var tests []models.Test
    err := r.db.Where("is_published = ?", true).Find(&tests).Error
    return tests, err
}

func (r *TestRepository) FindByType(testType string) ([]models.Test, error) {
    var tests []models.Test
    err := r.db.Where("type = ? AND is_published = ?", testType, true).Find(&tests).Error
    return tests, err
}

func (r *TestRepository) FindAll() ([]models.Test, error) {
    var tests []models.Test
    err := r.db.Find(&tests).Error
    return tests, err
}

func (r *TestRepository) Update(test *models.Test) error {
    return r.db.Save(test).Error
}

func (r *TestRepository) Delete(id string) error {
    return r.db.Delete(&models.Test{}, "id = ?", id).Error
}
