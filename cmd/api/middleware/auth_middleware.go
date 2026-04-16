package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gabrielleite03/kenjix_core/internal/service"
	"github.com/golang-jwt/jwt/v5"
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
		token = strings.Replace(token, "Bearer ", "", 1)
		/*
			valid, err := m.AuthService.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(w, "auth service error", http.StatusUnauthorized)
				return
			}

				if !valid {
					println("fodeu")
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}
		*/
		jwtService := NewJWTService()
		valid, err := jwtService.ValidateToken(token)
		if err != nil || !valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

type JWTService struct {
	secret []byte
}

func NewJWTService() *JWTService {
	// usa o MESMO secret do Java (não base64 aqui, igual você decidiu)
	secret := os.Getenv("REAL_JWT_SECRET")

	return &JWTService{
		secret: []byte(secret),
	}
}

// ValidateToken valida assinatura + expiração + algoritmo
func (s *JWTService) ValidateToken(tokenString string) (bool, error) {

	// remove "Bearer " se vier no header
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

		// garante HMAC (igual Java Algorithm.HMAC256)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return s.secret, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	// valida exp manual (opcional mas recomendado)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return false, fmt.Errorf("token expired")
		}
	}

	return true, nil
}
