package repository

import (
	"auth/internal/rest/models"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	CreateUser(user *models.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := ur.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := ur.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := ur.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := ur.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (ur *UserRepository) UpdateUser(user *models.User) error {
	if err := ur.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) DeleteUser(id uint) error {
	result := ur.db.Delete(&models.User{}, id)
	return result.Error
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	result := ur.db.Create(user)
	return result.Error
}
