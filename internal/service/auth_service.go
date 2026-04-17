package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AuthService interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
	Login(ctx context.Context, userName, password string) (*LoginResponse, error)
}

type authService struct {
	authURL string
	client  *http.Client
}

func NewAuthService(authURL string) AuthService {
	return &authService{
		authURL: authURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *authService) ValidateToken(ctx context.Context, token string) (bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		s.authURL+"/auth/validate",
		nil,
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", token)

	resp, err := s.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserName      string `json:"username"`
	Authenticated bool   `json:"authenticated"`
	Created       string `json:"created"`
	Expiration    string `json:"expiration"`
	AccessToken   string `json:"accessToken"`
	RefreshToken  string `json:"refreshToken"`
}

func (s *authService) Login(ctx context.Context, userName, password string) (*LoginResponse, error) {
	payload := LoginRequest{
		UserName: userName,
		Password: password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.authURL+"/auth/signin",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %s", string(bodyBytes))
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(bodyBytes, &loginResp); err != nil {
		return nil, err
	}

	return &loginResp, nil
}
