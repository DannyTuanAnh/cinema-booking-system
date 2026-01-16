package ticket_clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ticketHTTPClient struct {
	baseURL string
	http    *http.Client
}

func NewTicketHTTPClient(baseURL string) TicketClient {
	return &ticketHTTPClient{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *ticketHTTPClient) GetTicketByUserID(userID int64) ([]*GetTicketByUserIdResponse, error) {
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/tickets?user_id=%d", c.baseURL, userID), nil)
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {

		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("get ticket by user_id failed with status %d", resp.StatusCode)
		}

		if errMsg, ok := errorResp["error"]; ok {
			return nil, fmt.Errorf("%v", errMsg)
		}

		return nil, fmt.Errorf("get ticket by user_id with status %d", resp.StatusCode)
	}

	var ticketResp []*GetTicketByUserIdResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticketResp); err != nil {
		return nil, err
	}

	return ticketResp, nil
}
