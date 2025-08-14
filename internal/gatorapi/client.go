// Package gatorapi provides client functionality for interacting with APIs.
package gatorapi

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type GatorClient struct {
	client *http.Client 
}

func NewClient(timeout time.Duration) *GatorClient {
	return &GatorClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *GatorClient) Get(URL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to url %s: %w", URL, err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to url %s: %w", URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("received non-2xx status code %d from url %s", resp.StatusCode, URL)
    }

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}