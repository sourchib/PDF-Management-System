package repository

import (
	"database/sql"
	"pdf-management-system/internal/model"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (name, email, password, address, phone_number, post_code, role_id, is_email_verified, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	// Assuming default not verified, created_date is now
	return r.DB.QueryRow(query, user.Name, user.Email, user.Password, user.Address, user.PhoneNumber, user.PostCode, user.RoleID, user.IsEmailVerified, time.Now()).Scan(&user.ID)
}

func (r *UserRepository) FindUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, address, phone_number, post_code, role_id, is_email_verified, created_date
		FROM users
		WHERE email = $1
	`
	var user model.User
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Address, &user.PhoneNumber, &user.PostCode, &user.RoleID, &user.IsEmailVerified, &user.CreatedDate,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindById needed for middleware probably
func (r *UserRepository) FindUserByID(id int64) (*model.User, error) {
	query := `SELECT id, name, email, role_id FROM users WHERE id = $1`
	var user model.User
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.RoleID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
