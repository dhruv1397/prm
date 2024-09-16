package harness

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Get(ctx context.Context, client *http.Client, pat string, url string, responseDTO any) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error while forming request: %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("x-api-key", pat)

	response, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("error while executing request: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error while parsing response body: %w", err)
	}
	if response.StatusCode == http.StatusNotFound || len(body) == 0 {
		return nil
	}
	err = json.Unmarshal(body, responseDTO)
	if err != nil {
		return fmt.Errorf("error while parsing response: %w", err)
	}
	return nil
}
