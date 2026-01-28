package model

import "time"

type Role struct {
	ID   int64  `json:"id"`
	Role string `json:"role"`
}

type User struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Password        string     `json:"-"` // Don't return password in JSON
	Address         string     `json:"address"`
	PhoneNumber     string     `json:"phone_number"`
	PostCode        string     `json:"post_code"`
	RoleID          int64      `json:"role_id"`
	IsEmailVerified bool       `json:"is_email_verified"` // tinyint usually maps to bool or int8
	CreatedBy       *int64     `json:"created_by"`
	CreatedDate     time.Time  `json:"created_date"`
	ModifiedBy      *int64     `json:"modified_by"`
	ModifiedDate    *time.Time `json:"modified_date"`
}

type RegisterRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	PostCode    string `json:"post_code"`
	RoleID      int64  `json:"role_id"` // User sends role preference? Or default. Let's allow sending.
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
