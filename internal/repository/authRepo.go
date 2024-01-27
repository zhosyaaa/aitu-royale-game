package repository

import (
	"auth/internal/rest/models"
	"database/sql"
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
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := ur.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.Username, &user.Email, &user.Password, &user.Bank, &user.Awards,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.Username, &user.Email, &user.Password, &user.Bank, &user.Awards,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.Username, &user.Email, &user.Password, &user.Bank, &user.Awards,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetAllUsers() ([]models.User, error) {
	rows, err := ur.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
			&user.Username, &user.Email, &user.Password, &user.Bank, &user.Awards,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(user *models.User) error {
	_, err := ur.db.Exec(
		"UPDATE users SET username=$1, email=$2, password=$3, bank=$4, awards=$5 WHERE id=$6",
		user.Username, user.Email, user.Password, user.Bank, user.Awards, user.ID,
	)
	return err
}

func (ur *UserRepository) DeleteUser(id uint) error {
	_, err := ur.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	_, err := ur.db.Exec(
		"INSERT INTO users (created_at, updated_at, deleted_at, username, email, password, bank, awards) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.Username, user.Email, user.Password, user.Bank, user.Awards,
	)
	return err
}
