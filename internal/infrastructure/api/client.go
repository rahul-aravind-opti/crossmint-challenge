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
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := c.baseURL + endpoint

	// Create a retryable function
	retryableFunc := func(ctx context.Context) error {
		// Apply rate limiting
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return fmt.Errorf("rate limiter error: %w", err)
		}

		// Reset body reader for retry attempts
		if body != nil && bodyReader != nil {
			if seeker, ok := bodyReader.(io.Seeker); ok {
				_, err := seeker.Seek(0, io.SeekStart)
				if err != nil {
					return fmt.Errorf("failed to reset request body: %w", err)
				}
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		// Check if the response indicates an error
		if resp.StatusCode >= 500 {
			// Server errors are retryable
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return &temporaryError{
				error: domain.NewAPIError(resp.StatusCode, string(body), endpoint),
			}
		} else if resp.StatusCode >= 400 {
			// Client errors are not retryable
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			apiErr := domain.NewAPIError(resp.StatusCode, string(body), endpoint)

			if resp.StatusCode == http.StatusTooManyRequests {
				// Treat rate limit responses as retryable
				return &temporaryError{error: apiErr}
			}

			return &permanentError{error: apiErr}
		}

		// Success - store response for return
		return &successResponse{resp: resp}
	}

	// Execute with retry using our custom retry package
	err := pkgretry.Do(ctx, retryableFunc, c.retryConfig, isRetryableError)
	if err != nil {
		// Check if it's a success response
		var success *successResponse
		if errors.As(err, &success) {
			return success.resp, nil
		}
		return nil, err
	}

	return nil, fmt.Errorf("unexpected nil error from retry")
}

// temporaryError wraps errors that should be retried
type temporaryError struct {
	error
}

// permanentError wraps errors that should not be retried
type permanentError struct {
	error
}

// successResponse wraps a successful HTTP response
type successResponse struct {
	resp *http.Response
}

func (s *successResponse) Error() string {
	return "success"
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// successResponse is a special case - not really an error
	if _, ok := err.(*successResponse); ok {
		return false
	}

	// Check if it's a temporary error
	_, isTemporary := err.(*temporaryError)
	return isTemporary
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
