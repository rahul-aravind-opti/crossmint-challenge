package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	retry "github.com/avast/retry-go/v4"
	"github.com/crossmint/megaverse-challenge/internal/domain"
	"github.com/crossmint/megaverse-challenge/pkg/ratelimit"
	pkgretry "github.com/crossmint/megaverse-challenge/pkg/retry"
)

// Client represents the HTTP client for the Megaverse API
type Client struct {
	baseURL     string
	candidateID string
	httpClient  *http.Client
	rateLimiter *ratelimit.Limiter
	retryConfig pkgretry.Config
}

// ClientConfig holds the configuration for the API client
type ClientConfig struct {
	BaseURL           string
	CandidateID       string
	Timeout           time.Duration
	RetryConfig       pkgretry.Config
	RequestsPerSecond float64
}

// NewClient creates a new API client
func NewClient(config ClientConfig) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.RequestsPerSecond == 0 {
		config.RequestsPerSecond = 2.0
	}

	return &Client{
		baseURL:     config.BaseURL,
		candidateID: config.CandidateID,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		rateLimiter: ratelimit.NewLimiter(config.RequestsPerSecond),
		retryConfig: config.RetryConfig,
	}
}

// doRequest performs an HTTP request with rate limiting and retry logic
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	url := c.baseURL + endpoint
	var resp *http.Response

	retryableErr := pkgretry.Do(ctx, func(ctx context.Context) error {
		// Make every call synchronise on the limiter so bursts across goroutines keep a consistent pace.
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return retry.Unrecoverable(fmt.Errorf("rate limiter error: %w", err))
		}

		var bodyReader io.Reader
		if len(payload) > 0 {
			bodyReader = bytes.NewReader(payload)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return retry.Unrecoverable(fmt.Errorf("failed to create request: %w", err))
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		if resp != nil {
			// The same response pointer is reused by retry-go; close the previous body before issuing another attempt.
			resp.Body.Close()
			resp = nil
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			return err
		}

		status := resp.StatusCode
		if status == http.StatusTooManyRequests || status >= 500 {
			responseBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			resp = nil
			return domain.NewAPIError(status, string(responseBody), endpoint)
		}

		if status >= 400 {
			responseBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			resp = nil
			return retry.Unrecoverable(domain.NewAPIError(status, string(responseBody), endpoint))
		}

		return nil
	}, c.retryConfig, func(err error) bool {
		if err == nil {
			return false
		}

		// retry-go marks unrecoverable errors explicitly; respect that signal before applying our custom rules.
		if !retry.IsRecoverable(err) {
			return false
		}

		var apiErr *domain.APIError
		if errors.As(err, &apiErr) {
			return apiErr.StatusCode == http.StatusTooManyRequests || apiErr.StatusCode >= 500
		}

		return true
	})

	if retryableErr != nil {
		return nil, retryableErr
	}

	if resp == nil {
		return nil, fmt.Errorf("no response received from %s %s", method, endpoint)
	}

	return resp, nil
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, endpoint string, body interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return domain.NewAPIError(resp.StatusCode, string(body), endpoint)
	}

	return nil
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, endpoint string, body interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return domain.NewAPIError(resp.StatusCode, string(body), endpoint)
	}

	return nil
}

// Get performs a GET request and unmarshals the response
func (c *Client) Get(ctx context.Context, endpoint string, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return domain.NewAPIError(resp.StatusCode, string(body), endpoint)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// GetCandidateID returns the configured candidate ID
func (c *Client) GetCandidateID() string {
	return c.candidateID
}
