package book_clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type bookHTTPClient struct {
	baseURL string
	http    *http.Client
}

func NewBookHTTPClient(baseURL string) BookClient {
	return &bookHTTPClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *bookHTTPClient) BookSeats(req BookRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, c.baseURL+"/book", bytes.NewBuffer(body))
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("booking failed with status %d", resp.StatusCode)
		}
		if errMsg, ok := errorResp["error"]; ok {
			return fmt.Errorf("%v", errMsg)
		}
		return fmt.Errorf("booking failed with status %d", resp.StatusCode)
	}

	return nil

}
