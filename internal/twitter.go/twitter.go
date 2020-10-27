package twitter

import (
	"fmt"
	"net/http"
	"net/url"
)

const getRecentTweetsEndpoint = "https://api.twitter.com/2/tweets/search/recent"

// Client ...
type Client struct {
	HTTPClient *http.Client
	token      string
}

// Query ...
type Query struct {
	From        string
	TweetFields string // comma delimited string of fields you wish to search for
}

// New ...
func New(bearerToken string, client *http.Client) *Client {
	return &Client{
		HTTPClient: client,
		token:      bearerToken,
	}
}

// GetRecentTweets ...
func (c *Client) GetRecentTweets(query *Query) error {
	rawURL := fmt.Sprintf("%s?query=from:%s&twitter.fields=%s", getRecentTweetsEndpoint, query.From, query.TweetFields)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    parsedURL,
		Header: http.Header{
			"Authorization": []string{
				fmt.Sprintf("Bearer %s", c.token),
			},
		},
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not make http request: %s", resp.Status)
	}

	return nil

}
