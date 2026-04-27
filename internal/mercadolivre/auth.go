package mercadolivre

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type AuthService struct {
	ClientID     string
	ClientSecret string
	Token        *Token
}

func NewAuthService(clientID, clientSecret string) *AuthService {
	return &AuthService{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
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

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		UserID       int64  `json:"user_id"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	return &Token{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
	}, nil
}

func (a *AuthService) GetValidToken(ctx context.Context) (string, error) {
	if time.Now().Before(a.Token.ExpiresAt) {
		return a.Token.AccessToken, nil
	}

	t, err := a.RefreshToken(ctx)
	if err != nil {
		return "", err
	}

	a.Token = t
	return t.AccessToken, nil
}

func (a *AuthService) RefreshToken(ctx context.Context) (*Token, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", a.ClientID)
	form.Set("client_secret", a.ClientSecret)
	form.Set("refresh_token", a.Token.RefreshToken)

	req, _ := http.NewRequestWithContext(ctx, "POST",
		"https://api.mercadolibre.com/oauth/token",
		strings.NewReader(form.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	return &Token{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
	}, nil
}
