package movie_clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type movieHTTPClient struct {
	baseURL string
	http    *http.Client
}

func NewMovieHTTPClient(baseURL string) MovieClient {
	return &movieHTTPClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *movieHTTPClient) GetMovieDetails(ctx context.Context) ([]*GetMovieResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/movies", c.baseURL), nil)
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
			return nil, fmt.Errorf("get movie details failed with status %d", resp.StatusCode)
		}
		if errMsg, ok := errorResp["error"]; ok {
			return nil, fmt.Errorf("%v", errMsg)
		}
		return nil, fmt.Errorf("get movie details failed with status %d", resp.StatusCode)
	}

	var movieResp []*GetMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&movieResp); err != nil {
		return nil, err
	}

	return movieResp, nil
}
