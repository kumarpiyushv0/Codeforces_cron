package codeforces

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Contest represents a contest returned by the Codeforces API.
type Contest struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Phase               string `json:"phase"`
	Frozen              bool   `json:"frozen"`
	DurationSeconds     int64  `json:"durationSeconds"`
	StartTimeSeconds    int64  `json:"startTimeSeconds"`
	RelativeTimeSeconds int64  `json:"relativeTimeSeconds"`
}

// APIResponse wraps the standard Codeforces API response format.
type APIResponse struct {
	Status  string    `json:"status"`
	Comment string    `json:"comment,omitempty"`
	Result  []Contest `json:"result,omitempty"`
}

// Client interacts with the Codeforces API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Codeforces API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    "https://codeforces.com/api",
	}
}

// FetchUpcomingContests retrieves all contests with phase "BEFORE" (upcoming).
func (c *Client) FetchUpcomingContests(ctx context.Context) ([]Contest, error) {
	reqURL := fmt.Sprintf("%s/contest.list?gym=false", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create codeforces request: %w", err)
	}

	req.Header.Set("User-Agent", "CodeforcesCalendarSyncCron/1.0 (Go)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute codeforces request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("codeforces api returned status %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode codeforces response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("codeforces api error: %s", apiResp.Comment)
	}

	var upcoming []Contest
	for _, contest := range apiResp.Result {
		// phase "BEFORE" represents upcoming scheduled contests
		if contest.Phase == "BEFORE" && contest.StartTimeSeconds > 0 {
			upcoming = append(upcoming, contest)
		}
	}

	return upcoming, nil
}
