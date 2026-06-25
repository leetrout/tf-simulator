// Package client is a thin HTTP client for the tfsim cloud REST API
// exposed at /api/cloud/*. It uses only the Go standard library.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client talks to a running tfsim server.
type Client struct {
	Endpoint   string
	HTTPClient *http.Client
}

// New returns a Client pointed at endpoint (e.g. "http://localhost:9321").
func New(endpoint string) *Client {
	return &Client{
		Endpoint:   strings.TrimRight(endpoint, "/"),
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NotFoundError is returned when the server responds with 404.
type NotFoundError struct {
	Collection string
	ID         string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s/%s not found", e.Collection, e.ID)
}

// IsNotFound reports whether err is a NotFoundError.
func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func (c *Client) url(collection, id string) string {
	base := fmt.Sprintf("%s/api/cloud/%s", c.Endpoint, collection)
	if id == "" {
		return base
	}
	return base + "/" + id
}

// do performs a request and decodes a JSON object response into out (if non-nil).
func (c *Client) do(method, url string, body any, out any, collection, id string) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{Collection: collection, ID: id}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s: unexpected status %d: %s", method, url, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

// Create POSTs a new object into the collection and decodes the result.
func (c *Client) Create(collection string, body any, out any) error {
	return c.do(http.MethodPost, c.url(collection, ""), body, out, collection, "")
}

// Get fetches an object by id. Returns a NotFoundError on 404.
func (c *Client) Get(collection, id string, out any) error {
	return c.do(http.MethodGet, c.url(collection, id), nil, out, collection, id)
}

// Update PATCHes an object by id (partial update) and decodes the result.
func (c *Client) Update(collection, id string, body any, out any) error {
	return c.do(http.MethodPatch, c.url(collection, id), body, out, collection, id)
}

// Delete removes an object by id.
func (c *Client) Delete(collection, id string) error {
	return c.do(http.MethodDelete, c.url(collection, id), nil, nil, collection, id)
}
