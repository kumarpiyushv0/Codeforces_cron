package codeforces

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchUpcomingContests(t *testing.T) {
	mockResponse := `{
		"status": "OK",
		"result": [
			{
				"id": 2001,
				"name": "Codeforces Round 1000 (Div. 2)",
				"type": "CF",
				"phase": "BEFORE",
				"frozen": false,
				"durationSeconds": 7200,
				"startTimeSeconds": 1729000000,
				"relativeTimeSeconds": -86400
			},
			{
				"id": 2000,
				"name": "Codeforces Round 999 (Div. 1)",
				"type": "CF",
				"phase": "FINISHED",
				"frozen": false,
				"durationSeconds": 7200,
				"startTimeSeconds": 1728000000,
				"relativeTimeSeconds": 86400
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	contests, err := client.FetchUpcomingContests(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(contests) != 1 {
		t.Fatalf("expected 1 upcoming contest, got %d", len(contests))
	}

	if contests[0].ID != 2001 {
		t.Errorf("expected contest ID 2001, got %d", contests[0].ID)
	}
	if contests[0].Name != "Codeforces Round 1000 (Div. 2)" {
		t.Errorf("expected name 'Codeforces Round 1000 (Div. 2)', got %s", contests[0].Name)
	}
}
