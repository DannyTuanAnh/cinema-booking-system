package show_clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type showHTTPClient struct {
	baseURL string
	http    *http.Client
}

func NewShowHTTPClient(baseURL string) ShowClient {
	return &showHTTPClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *showHTTPClient) GetShowByMovieID(movieID int) ([]*GetShowResponse, error) {
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/shows?movie_id=%d", c.baseURL, movieID), nil)
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
			return nil, fmt.Errorf("get show by movie ID failed with status %d", resp.StatusCode)
		}
		if errMsg, ok := errorResp["error"]; ok {
			return nil, fmt.Errorf("%v", errMsg)
		}
		return nil, fmt.Errorf("get show by movie ID failed with status %d", resp.StatusCode)
	}

	var showResp []*GetShowResponse
	if err := json.NewDecoder(resp.Body).Decode(&showResp); err != nil {
		return nil, err
	}

	return showResp, nil
}
