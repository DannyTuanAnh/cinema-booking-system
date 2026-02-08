package auth_clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type authHTTPClient struct {
	baseURL string
	http    *http.Client
}

func NewAuthHTTPClient(baseURL string) AuthClient {
	return &authHTTPClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *authHTTPClient) Login(ctx context.Context, req LoginRequest) (*LoginResponse, []*http.Cookie, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/auth/login", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Kiểm tra status code từ core server
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, nil, fmt.Errorf("login failed with status %d", resp.StatusCode)
		}
		if errMsg, ok := errorResp["error"]; ok {
			return nil, nil, fmt.Errorf("%v", errMsg)
		}
		return nil, nil, fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, nil, err
	}

	return &loginResp, resp.Cookies(), nil
}

func (c *authHTTPClient) Register(ctx context.Context, req RegisterRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/auth/register", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Kiểm tra status code từ core server
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("registration failed with status %d", resp.StatusCode)
		}
		if errMsg, ok := errorResp["error"]; ok {
			return fmt.Errorf("%v", errMsg)
		}
		return fmt.Errorf("registration failed with status %d", resp.StatusCode)
	}

	return nil
}

func (c *authHTTPClient) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, []*http.Cookie, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/auth/refresh", nil)
	if err != nil {
		return nil, nil, err
	}

	httpReq.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	// Kiểm tra status code từ core server
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, nil, fmt.Errorf("refresh failed with status %d", resp.StatusCode)
		}

		if errMsg, ok := errorResp["error"]; ok {
			return nil, nil, fmt.Errorf("%v", errMsg)
		}

		return nil, nil, fmt.Errorf("refresh failed with status %d", resp.StatusCode)
	}

	var refreshResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return nil, nil, err
	}

	return &refreshResp, resp.Cookies(), nil
}
