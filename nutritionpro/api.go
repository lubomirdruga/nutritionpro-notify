package nutritionpro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseUrl = "https://api.nutritionpro.eu"

// Client represents the API client configuration
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

// GetMenu retrieves the menu from the API
func (c *Client) GetMenu() (*MenuResponse, error) {
	req, err := c.newRequest(http.MethodGet, "/api/menu/me", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var menu MenuResponse
	if err := c.doRequest(req, &menu); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return &menu, nil
}

// newRequest creates a new HTTP request with proper headers
func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// doRequest executes the HTTP request and unmarshalls the response
func (c *Client) doRequest(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("unmarshaling response: %w", err)
	}

	return nil
}

// LoginRequest represents the authentication request payload
type LoginRequest struct {
	PhoneNumber string `json:"inBodyId"`
}

// LoginResponse represents the authentication response
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

// GetToken authenticates with the API and returns the access token
func (c *Client) GetToken(phoneNumber string) (*LoginResponse, error) {
	loginData := LoginRequest{
		PhoneNumber: phoneNumber,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return nil, fmt.Errorf("marshaling login data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/api/menu/rate/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("creating login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	var loginResponse LoginResponse
	if err := c.doRequest(req, &loginResponse); err != nil {
		return nil, fmt.Errorf("executing login request: %w", err)
	}

	c.apiToken = loginResponse.AccessToken

	return &loginResponse, nil
}

// NewClientWithAuth creates a new client and authenticates immediately
func NewClientWithAuth(phoneNumber string) (*Client, error) {
	client := &Client{
		baseURL: baseUrl,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Authenticate and get a token
	loginResp, err := client.GetToken(phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("authenticating client: %w", err)
	}

	client.apiToken = loginResp.AccessToken
	return client, nil
}
