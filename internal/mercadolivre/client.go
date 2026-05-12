package mercadolivre

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
	Auth    *AuthService
	UserID  int64
}

func (c *Client) Get(ctx context.Context, path string, out interface{}) error {
	token, err := c.Auth.GetValidToken(ctx, c.UserID)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("ML error: %s", string(body))
	}

	return json.Unmarshal(body, out)
}
