package reddit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"rtrade/config"
)

const (
    tokenURL = "https://www.reddit.com/api/v1/access_token"
	userAgent = "web:RedditTrader:v0.0.1 (by /u/cougargriff)"
	apiURL   = "https://oauth.reddit.com"
)


// TokenResponse represents the OAuth token response from Reddit
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    Scope        string `json:"scope"`
    RefreshToken string `json:"refresh_token,omitempty"`
}

// Client handles Reddit OAuth operations
type Client struct {
    config     * config.RedditConfig
    httpClient *http.Client
}

// NewClient creates a new Reddit OAuth client
func NewClient(config * config.RedditConfig) *Client {
    return &Client{
        config: config,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// ExchangeCode exchanges an authorization code for an access token
func (c *Client) ExchangeCode(code string) (*TokenResponse, error) {
    // Prepare form data
    data := url.Values{}
    data.Set("grant_type", "authorization_code")
    data.Set("code", code)
    data.Set("redirect_uri", c.config.RedirectUrl)

    // Create request
    req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", userAgent)

    // Set Basic Auth header
    auth := base64.StdEncoding.EncodeToString([]byte(c.config.ClientId + ":" + c.config.ClientSecret))
    req.Header.Set("Authorization", "Basic "+auth)

    // Execute request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Check status code
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    // Parse response
    var tokenResp TokenResponse
    if err := json.Unmarshal(body, &tokenResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    return &tokenResp, nil
}

// RefreshToken refreshes an access token using a refresh token
func (c *Client) RefreshToken(refreshToken string) (*TokenResponse, error) {
    // Prepare form data
    data := url.Values{}
    data.Set("grant_type", "refresh_token")
    data.Set("refresh_token", refreshToken)

    // Create request
    req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", userAgent)

    // Set Basic Auth header
    auth := base64.StdEncoding.EncodeToString([]byte(c.config.ClientId + ":" + c.config.ClientSecret))
    req.Header.Set("Authorization", "Basic "+auth)

    // Execute request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Check status code
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    // Parse response
    var tokenResp TokenResponse
    if err := json.Unmarshal(body, &tokenResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    return &tokenResp, nil
}

