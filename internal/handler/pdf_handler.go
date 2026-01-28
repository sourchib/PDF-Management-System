package handler

import (
	"encoding/json"
	"net/http"
	"pdf-management-system/internal/model"
	"pdf-management-system/internal/service"
	"strconv"
	"strings"
)

type PdfHandler struct {
	Service *service.PdfService
}

func NewPdfHandler(service *service.PdfService) *PdfHandler {
	return &PdfHandler{Service: service}
}

func (h *PdfHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.GeneratePdfRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body", "")
		return
	}

	pdf, err := h.Service.GeneratePDF(req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}

	respondSuccess(w, "PDF generated successfully", pdf)
}

func (h *PdfHandler) UploadPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 10MB limit
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "File size exceeds maximum limit (10MB)", "FILE_TOO_LARGE")
		return
	}

	file, header, err := r.FormFile("file") // Field name wasn't specified but usually 'file'
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing file part", "")
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/pdf" && !strings.HasSuffix(header.Filename, ".pdf") {
		respondError(w, http.StatusBadRequest, "Only PDF files are allowed", "INVALID_FILE_TYPE")
		return
	}

	pdf, err := h.Service.UploadPDF(file, header)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}

	respondSuccess(w, "PDF uploaded successfully", pdf)
}

func (h *PdfHandler) ListPDFs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	files, total, err := h.Service.ListPDFs(status, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	resp := model.PaginatedResponse{
		Success: true,
		Data:    files,
		Pagination: model.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *PdfHandler) DeletePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path: /api/pdf/{id}
	// We assume registered as /api/pdf/
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 0 {
		respondError(w, http.StatusBadRequest, "Invalid URL", "")
		return
	}
	idStr := parts[len(parts)-1]
	// Handle trailing slash
	if idStr == "" && len(parts) > 1 {
		idStr = parts[len(parts)-2]
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID", "")
		return
	}

	pdf, err := h.Service.DeletePDF(id)
	if err != nil {
		// Distinguish not found vs other errors?
		if err.Error() == "file not found" {
			respondError(w, http.StatusNotFound, "File not found", "")
		} else if err.Error() == "file already deleted" {
			respondError(w, http.StatusBadRequest, "File already deleted", "")
		} else {
			respondError(w, http.StatusInternalServerError, err.Error(), "")
		}
		return
	}

	respondSuccess(w, "PDF deleted successfully", pdf)
}

func respondSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func respondError(w http.ResponseWriter, code int, message string, errorCode string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(model.ApiResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
	})
}
