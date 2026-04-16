package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type AuthHandler struct {
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// POST /login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	resp, err := h.AuthService.Login(
		r.Context(),
		req.UserName,
		req.Password,
	)

	if err != nil {
		print(err.Error())
		http.Error(w, "auth service error", http.StatusInternalServerError)
		return
	}

	if resp == nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
