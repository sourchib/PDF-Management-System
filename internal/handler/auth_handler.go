package handler

import (
	"encoding/json"
	"net/http"
	"pdf-management-system/internal/model"
	"pdf-management-system/internal/service"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body", "")
		return
	}

	user, err := h.Service.Register(req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error(), "REGISTRATION_FAILED")
		return
	}

	respondSuccess(w, "User registered successfully", user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body", "")
		return
	}

	resp, err := h.Service.Login(req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error(), "LOGIN_FAILED")
		return
	}

	respondSuccess(w, "Login successful", resp)
}
