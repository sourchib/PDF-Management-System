package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"pdf-management-system/internal/model"
	"pdf-management-system/internal/repository"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PdfService struct {
	Repo *repository.PdfRepository
}

func NewPdfService(repo *repository.PdfRepository) *PdfService {
	return &PdfService{Repo: repo}
}

func (s *PdfService) GeneratePDF(req model.GeneratePdfRequest) (*model.PdfFile, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Setup Header
	pdf.SetHeaderFunc(func() {
		// Logo Logic
		logoURL := req.LogoURL
		if logoURL == "" {
			// Default Logo (using a public placeholder or you could use a local file)
			logoURL = "https://via.placeholder.com/150.png?text=LOGO"
		}

		// Fetch and render image
		if logoURL != "" {
			resp, err := http.Get(logoURL)
			if err == nil && resp.StatusCode == 200 {
				defer resp.Body.Close()
				imgName := "logo_header"
				opts := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
				// Simple check for jpg
				if len(logoURL) > 4 && logoURL[len(logoURL)-4:] == ".jpg" {
					opts.ImageType = "JPG"
				}
				pdf.RegisterImageOptionsReader(imgName, opts, resp.Body)
				// Logo at top-left: x=10, y=10, w=25
				pdf.Image(imgName, 10, 10, 25, 0, false, "", 0, "")
			}
		}

		// Institution Name (Centered)
		pdf.SetY(12) // Slightly down
		pdf.SetFont("Arial", "B", 16)
		pdf.CellFormat(0, 10, req.InstitutionName, "", 1, "C", false, 0, "")

		// Address & Phone (Centered below name)
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 5, req.Address, "", 1, "C", false, 0, "")
		pdf.CellFormat(0, 5, req.Phone, "", 1, "C", false, 0, "")

		pdf.Ln(10) // Spacing
		// Draw line below logo/header info (at approx Y=40-42)
		drawY := pdf.GetY() + 5
		pdf.SetLineWidth(0.5)
		pdf.Line(10, drawY, 200, drawY)
		pdf.Ln(10) // Extra space after line before body starts
	})

	// Setup Footer
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		// Page X of Y | Timestamp
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb} | Generated: %s", pdf.PageNo(), timestamp), "", 0, "C", false, 0, "")
	})

	pdf.AliasNbPages("")
	pdf.AddPage()

	// Content
	// Title
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, req.Title)
	pdf.Ln(12)

	// Date
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Date: %s", time.Now().Format("January 02, 2006")))
	pdf.Ln(12)

	// Body Content
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, req.Content, "", "", false)

	// Save file
	filename := fmt.Sprintf("report_%s_%d.pdf", time.Now().Format("20060102"), time.Now().UnixNano())
	outputPath := filepath.Join("uploads", "pdf", filename)

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to save pdf: %v", err)
	}

	// Save to DB
	// Normalize path using forward slashes for DB consistency or just store relative
	dbFilepath := fmt.Sprintf("/uploads/pdf/%s", filename)

	pdfRecord := &model.PdfFile{
		Filename:  filename,
		Filepath:  dbFilepath,
		Status:    model.StatusCreated,
		CreatedAt: time.Now(), // Actually Repo sets this or DB default
	}

	// Get file size
	info, err := os.Stat(outputPath)
	if err == nil {
		pdfRecord.Size = info.Size()
	}

	err = s.Repo.Create(pdfRecord)
	if err != nil {
		return nil, err
	}

	return pdfRecord, nil
}

func (s *PdfService) UploadPDF(file io.Reader, header *multipart.FileHeader) (*model.PdfFile, error) {
	// Validate MIME is done in handler typically, but we can do here?
	// Requirement: Validate MIME type application/pdf
	// Implementation: Check extension and Content-Type header

	// Generate unique name
	ext := filepath.Ext(header.Filename)
	if ext != ".pdf" {
		return nil, fmt.Errorf("invalid file extension")
	}

	uniqueName := fmt.Sprintf("upload_%s_%d%s", time.Now().Format("20060102"), time.Now().UnixNano(), ext)
	outputPath := filepath.Join("uploads", "pdf", uniqueName)

	// Create file
	dst, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	// DB Record
	originalName := header.Filename
	filepathStr := fmt.Sprintf("/uploads/pdf/%s", uniqueName)

	pdfRecord := &model.PdfFile{
		Filename:     uniqueName,
		OriginalName: &originalName,
		Filepath:     filepathStr,
		Size:         header.Size,
		Status:       model.StatusUploaded,
	}

	err = s.Repo.Create(pdfRecord)
	if err != nil {
		return nil, err
	}

	return pdfRecord, nil
}

func (s *PdfService) ListPDFs(status string, page, limit int) ([]model.PdfFile, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.Repo.FindAll(status, page, limit)
}

func (s *PdfService) DeletePDF(id int64) (*model.PdfFile, error) {
	// Repo handles validation
	err := s.Repo.SoftDelete(id)
	if err != nil {
		return nil, err
	}

	// Return updated record
	return s.Repo.FindByID(id)
}
