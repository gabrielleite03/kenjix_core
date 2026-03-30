package middleware

import (
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type AuthMiddleware struct {
	AuthService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
	}
}

func (m *AuthMiddleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		valid, err := m.AuthService.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "auth service error", http.StatusUnauthorized)
			return
		}

		if !valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
