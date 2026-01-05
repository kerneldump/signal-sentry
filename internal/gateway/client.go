package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"tmobile-stats/internal/models"
)

// FetchStats retrieves all gateway statistics from the T-Mobile Home Internet Gateway.
// It implements a basic retry mechanism for transient network or server errors.
func FetchStats(client *http.Client, url string) (*models.GatewayResponse, error) {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			// Basic backoff
			time.Sleep(time.Duration(i) * 100 * time.Millisecond)
		}

		data, err := fetchOnce(client, url)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("failed to fetch stats after %d attempts: %w", maxRetries, lastErr)
}

func fetchOnce(client *http.Client, url string) (*models.GatewayResponse, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data models.GatewayResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

