package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const perPage = 100

func GetUserActivity(username string) ([]Event, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=%d", username, perPage)

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "gh-user-activity/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	LogApiCall(string(body), resp.StatusCode, username)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return events, nil
}
