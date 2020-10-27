package twitter

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"os"
	"encoding/json"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"

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

// TweetResponse ...
type TweetResponse struct {
	Data []Tweet
}

// Tweet ...
type Tweet struct {
	Text string `json:"text,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	ID string `json:"id"`
}

// New ...
func New(bearerToken string, client *http.Client) *Client {
	return &Client{
		HTTPClient: client,
		token:      bearerToken,
	}
}

// GetRecentTweets ...
func (c *Client) GetRecentTweets(query *Query) (*TweetResponse, error) {
	rawURL := fmt.Sprintf("%s?query=from:%s&twitter.fields=%s", getRecentTweetsEndpoint, query.From, query.TweetFields)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not make http request: %s", resp.Status)
	}
	var tweetResp TweetResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &tweetResp); err != nil {
		return nil, err
	}

	return &tweetResp, nil
}

// Handler ...
func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	client := New(os.Getenv("TWITTER_BEARER_TOKEN"), &http.Client{})
	if m.ChannelID == os.Getenv("TWEET_CHANNEL") && strings.Contains(m.Content, "dril") {
		tweets, err := client.GetRecentTweets(&Query{
			From: "@dril",
			TweetFields: "created_at,entities",
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, tweet := range tweets.Data {
			fmt.Println(tweet)
		}
	}

}
