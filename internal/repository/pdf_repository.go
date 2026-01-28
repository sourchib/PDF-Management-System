package repository

import (
	"database/sql"
	"fmt"
	"pdf-management-system/internal/model"
	"time"
)

type PdfRepository struct {
	DB *sql.DB
}

func NewPdfRepository(db *sql.DB) *PdfRepository {
	return &PdfRepository{DB: db}
}

func (r *PdfRepository) Create(pdf *model.PdfFile) error {
	query := `
		INSERT INTO pdf_files (filename, original_name, filepath, size, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.DB.QueryRow(query, pdf.Filename, pdf.OriginalName, pdf.Filepath, pdf.Size, pdf.Status, time.Now()).Scan(&pdf.ID)
}

func (r *PdfRepository) FindByID(id int64) (*model.PdfFile, error) {
	query := `
		SELECT id, filename, original_name, filepath, size, status, created_at, updated_at, deleted_at
		FROM pdf_files
		WHERE id = $1 AND status != 'DELETED'
	`
	var pdf model.PdfFile
	err := r.DB.QueryRow(query, id).Scan(
		&pdf.ID, &pdf.Filename, &pdf.OriginalName, &pdf.Filepath, &pdf.Size, &pdf.Status, &pdf.CreatedAt, &pdf.UpdatedAt, &pdf.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &pdf, nil
}

func (r *PdfRepository) FindAll(status string, page, limit int) ([]model.PdfFile, int64, error) {
	offset := (page - 1) * limit

	// Base query
	query := `SELECT id, filename, original_name, filepath, size, status, created_at, updated_at, deleted_at FROM pdf_files WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM pdf_files WHERE 1=1`
	args := []interface{}{}
	argId := 1

	if status != "" {
		filter := fmt.Sprintf(" AND status = $%d", argId)
		query += filter
		countQuery += filter
		args = append(args, status)
		argId++
	}

	// Always filter out deleted by default? PROMPT says "status yang jelas... DELETED - file yang sudah dihapus".
	// Requirement 3: "Support filter berdasarkan status".
	// Requirement 2: "Tampilkan semua file PDF yang ada di database... (deleted inclusive?)".
	// Usually "List" implies active files unless "DELETED" is asked.
	// But requirement lists "DELETED" as a status. "Setiap file harus memiliki status yang jelas... DELETED".
	// So we should probably list them if not filtered out, or maybe just list everything.
	// However, Task 4 says "Soft Delete... return error jika file sudah dalam status DELETED" (implied for delete action).
	// For List, I will list ALL including DELETED unless filtered.

	// Order by latest
	query += " ORDER BY created_at DESC"

	// Pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, limit, offset)

	// Execute count
	var total int64
	// We need args for count query (only status)
	countArgs := args[:argId-1]
	err := r.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var files []model.PdfFile
	for rows.Next() {
		var pdf model.PdfFile
		if err := rows.Scan(&pdf.ID, &pdf.Filename, &pdf.OriginalName, &pdf.Filepath, &pdf.Size, &pdf.Status, &pdf.CreatedAt, &pdf.UpdatedAt, &pdf.DeletedAt); err != nil {
			return nil, 0, err
		}
		files = append(files, pdf)
	}

	return files, total, nil
}

func (r *PdfRepository) SoftDelete(id int64) error {
	// Check if exists first? Or just update.
	// Requirement: Return error if file not found OR already deleted.

	// We can check first
	pdf, err := r.FindByID(id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("file not found")
	} else if err != nil {
		return err
	}

	if pdf.Status == model.StatusDeleted {
		return fmt.Errorf("file already deleted")
	}

	query := `
		UPDATE pdf_files
		SET status = 'DELETED', deleted_at = $1
		WHERE id = $2
	`
	_, err = r.DB.Exec(query, time.Now(), id)
	return err
}
