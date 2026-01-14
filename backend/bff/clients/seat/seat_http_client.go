package seat_clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type seatHTTPClient struct {
	baseURL string
	http    *http.Client
}

func NewSeatHTTPClient(baseURL string) SeatClient {
	return &seatHTTPClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *seatHTTPClient) GetSeatByShowID(showID int) ([]*GetSeatResponse, error) {
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/seats?show_id=%d", c.baseURL, showID), nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Kiểm tra status code từ core server
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("get seat by show ID failed with status %d", resp.StatusCode)
		}
		if errMsg, ok := errorResp["error"]; ok {
			return nil, fmt.Errorf("%v", errMsg)
		}
		return nil, fmt.Errorf("get seat by show ID failed with status %d", resp.StatusCode)
	}

	var seatResp []*GetSeatResponse
	if err := json.NewDecoder(resp.Body).Decode(&seatResp); err != nil {
		return nil, err
	}

	return seatResp, nil
}
