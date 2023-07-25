package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	grantType           = "client_credentials"
	accessTokenEndpoint = "https://id.twitch.tv/oauth2/token"
)

type Client struct {
	HttpClient   *http.Client
	ClientID     string
	ClientSecret string
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"tokenType"`
}

func (t *Client) Auth() (*AuthResponse, error) {
	return t.sendAuthRequest()
}

func (t *Client) sendAuthRequest() (*AuthResponse, error) {

	params := url.Values{}
	params.Set("client_id", t.ClientID)
	params.Set("client_secret", t.ClientSecret)
	params.Set("grant_type", grantType)

	parsedURL, err := url.Parse(accessTokenEndpoint)
	if err != nil {
		return nil, fmt.Errorf("twitch API Auth URL is malformed: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, parsedURL.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating new HTTP request for auth to Twitch API: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request for auth to Twitch API: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("response from HTTP request for auth to Twitch API is empty")
	} else {
		var authResponse AuthResponse
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body from Twitch API auth endpoint: %w", err)
		}
		if err := json.Unmarshal(body, &authResponse); err != nil {
			return nil, fmt.Errorf("error unmarshaling response body from Twitch API auth endpoint: %w", err)
		} else {
			return &authResponse, nil
		}
	}
}
