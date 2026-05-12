package mercadolivre

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/repository"
)

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       int64     `json:"user_id"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
}

type AuthService struct {
	ClientID     string
	ClientSecret string
	TokenDAO     repository.MLTokensDAO
}

func NewAuthService(clientID, clientSecret string, dao repository.MLTokensDAO) *AuthService {
	return &AuthService{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenDAO:     dao,
	}
}

func RedirectToMercadoLivreAuth(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("ML_CLIENT_ID")
	redirectURI := "https://kenjipet.com.br/redirect_mercadolivre"

	url := fmt.Sprintf(
		"https://auth.mercadolivre.com.br/authorization?response_type=code&client_id=%s&redirect_uri=%s",
		clientID,
		url.QueryEscape(redirectURI),
	)

	http.Redirect(w, r, url, http.StatusFound)
}

func GetToken(ctx context.Context, clientID, clientSecret, code, redirectURI string) (*Token, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, _ := http.NewRequestWithContext(ctx,
		"POST",
		"https://api.mercadolibre.com/oauth/token",
		strings.NewReader(form.Encode()),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("mercadolivre status=%d body=%s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
		UserID       int64  `json:"user_id"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar body: %w | body=%s", err, string(body))
	}

	return &Token{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
		UserID:       result.UserID,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (a *AuthService) GetValidToken(ctx context.Context, userID int64) (string, error) {

	token, err := a.TokenDAO.GetByUserID(userID)
	if err != nil {
		return "", err
	}

	if token == nil {
		return "", fmt.Errorf("token não encontrado")
	}

	// 🔥 token ainda válido
	if time.Now().Before(token.ExpiresAt) {
		return token.AccessToken, nil
	}

	// 🔥 token expirado → refresh
	newToken, err := a.RefreshToken(ctx, token.RefreshToken)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)

	err = a.TokenDAO.UpdateTokens(
		userID,
		newToken.AccessToken,
		newToken.RefreshToken,
		expiresAt,
	)
	if err != nil {
		return "", err
	}

	return newToken.AccessToken, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*Token, error) {

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", a.ClientID)
	form.Set("client_secret", a.ClientSecret)
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.mercadolibre.com/oauth/token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("erro refresh ML: %s", string(body))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (a *AuthService) SaveToken(ctx context.Context, token *Token) error {

	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	dbToken := &model.MLToken{
		UserID:       token.UserID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    expiresAt,
	}

	return a.TokenDAO.Upsert(dbToken)
}
