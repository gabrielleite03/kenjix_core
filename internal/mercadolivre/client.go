package mercadolivre

import (
	"context"
	"encoding/json"
	"net/http"
)

type Client struct {
	BaseURL string
	Auth    *AuthService
}

func (c *Client) Get(ctx context.Context, path string, out interface{}) error {
	token, err := c.Auth.GetValidToken(ctx)
	if err != nil {
		return err
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", c.BaseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}
