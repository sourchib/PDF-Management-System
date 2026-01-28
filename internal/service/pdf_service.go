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

	// Setup Footer (Recurring on all pages)
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		// Page X of Y | Timestamp
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb} | Generated: %s", pdf.PageNo(), timestamp), "", 0, "C", false, 0, "")
	})

	pdf.AliasNbPages("")
	pdf.AddPage()

	// --- HEADER PDF (Page 1 Only) ---
	// 1. Logo/image di kiri atas
	logoURL := req.LogoURL
	if logoURL == "" {
		// Default Logo
		logoURL = "https://via.placeholder.com/150.png?text=LOGO"
	}

	// Fetch logo
	resp, err := http.Get(logoURL)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		imgName := "logo_header"
		opts := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		if filepath.Ext(logoURL) == ".jpg" || filepath.Ext(logoURL) == ".jpeg" {
			opts.ImageType = "JPG"
		}
		pdf.RegisterImageOptionsReader(imgName, opts, resp.Body)
		// Logo di kiri atas
		pdf.Image(imgName, 10, 10, 25, 0, false, "", 0, "")
	}

	// 2. Nama institusi/perusahaan di tengah
	pdf.SetY(15)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, req.InstitutionName, "", 1, "C", false, 0, "")

	// 3. Alamat dan kontak di bawah nama
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, req.Address, "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Kontak: %s", req.Phone), "", 1, "C", false, 0, "")

	// Separator Line
	pdf.Ln(5)
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(10)

	// --- CONTENT PDF ---
	// Judul dokumen (dari parameter)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, req.Title)
	pdf.Ln(10)

	// Tanggal generate
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Tanggal Generate: %s", time.Now().Format("02 January 2006")))
	pdf.Ln(12)

	// Isi konten dokumen (support text/paragraf)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 8, req.Content, "", "L", false)

	// Save file
	filename := fmt.Sprintf("report_%s_%d.pdf", time.Now().Format("20060102"), time.Now().UnixNano())
	outputPath := filepath.Join("uploads", "pdf", filename)

	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to save pdf: %v", err)
	}

	// Save to DB
	dbFilepath := fmt.Sprintf("/uploads/pdf/%s", filename)
	pdfRecord := &model.PdfFile{
		Filename:  filename,
		Filepath:  dbFilepath,
		Status:    model.StatusCreated,
		CreatedAt: time.Now(),
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
	ext := filepath.Ext(header.Filename)
	if ext != ".pdf" {
		return nil, fmt.Errorf("hanya menerima file dengan ekstensi .pdf")
	}

	uniqueName := fmt.Sprintf("upload_%s_%d%s", time.Now().Format("20060102"), time.Now().UnixNano(), ext)
	outputPath := filepath.Join("uploads", "pdf", uniqueName)

	dst, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

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
	err := s.Repo.SoftDelete(id)
	if err != nil {
		return nil, err
	}
	return s.Repo.FindByID(id)
}
