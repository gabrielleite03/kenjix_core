package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gabrielleite03/kenjix_core/internal/config"
	"github.com/gabrielleite03/kenjix_core/internal/mercadolivre"
)

type MercadolivreHandler struct {
	service *mercadolivre.AuthService
}

func NewMercadolivreHandler(service *mercadolivre.AuthService) *MercadolivreHandler {
	return &MercadolivreHandler{service: service}
}

func (h *MercadolivreHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := mercadolivre.GetToken(
		r.Context(),
		config.Load().ClientID,
		config.Load().ClientSecret,
		code,
		"http://localhost:8080/callback",
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println("ACCESS TOKEN:", token.AccessToken)
	fmt.Println("REFRESH TOKEN:", token.RefreshToken)

	w.Write([]byte("Token gerado com sucesso!"))
}

func saveEnv(key, value string) error {
	file, err := os.OpenFile(".env", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	content, _ := io.ReadAll(file)

	lines := strings.Split(string(content), "\n")
	found := false

	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + value
			found = true
		}
	}

	if !found {
		lines = append(lines, key+"="+value)
	}

	newContent := strings.Join(lines, "\n")

	return os.WriteFile(".env", []byte(newContent), 0644)
}
