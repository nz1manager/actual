package repository

import (
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "yourmodule/internal/models"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
    var user models.User
    err := r.db.Where("id = ?", id).First(&user).Error
    return &user, err
}

func (r *UserRepository) FindByGoogleID(googleID string) (*models.User, error) {
    var user models.User
    err := r.db.Where("google_id = ?", googleID).First(&user).Error
    return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
    return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) FindAll() ([]models.User, error) {
    var users []models.User
    err := r.db.Where("is_admin = ?", false).Find(&users).Error
    return users, err
}
