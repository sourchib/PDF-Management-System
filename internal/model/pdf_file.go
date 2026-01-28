package model

import "time"

type PdfStatus string

const (
	StatusCreated  PdfStatus = "CREATED"
	StatusUploaded PdfStatus = "UPLOADED"
	StatusDeleted  PdfStatus = "DELETED"
)

type PdfFile struct {
	ID           int64      `json:"id"`
	Filename     string     `json:"filename"`
	OriginalName *string    `json:"original_name"`
	Filepath     string     `json:"filepath"`
	Size         int64      `json:"size"`
	Status       PdfStatus  `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type GeneratePdfRequest struct {
	Title           string `json:"title"`
	InstitutionName string `json:"institution_name"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	LogoURL         string `json:"logo_url"`
	Content         string `json:"content"`
}

type ApiResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
}

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}
