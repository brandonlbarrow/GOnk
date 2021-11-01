package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

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
	Text      string `json:"text,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	ID        string `json:"id"`
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
	rawURL := fmt.Sprintf("%s?query=from:%s&tweet.fields=%s", getRecentTweetsEndpoint, query.From, query.TweetFields)
	fmt.Println(rawURL)
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
	user := os.Getenv("TWEET_USER")
	text := os.Getenv("TWEET_LISTEN_TEXT")

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	client := New(os.Getenv("TWITTER_BEARER_TOKEN"), &http.Client{})
	if m.ChannelID == os.Getenv("TWEET_CHANNEL") && strings.Contains(m.Content, text) {
		tweets, err := client.GetRecentTweets(&Query{
			From:        user,
			TweetFields: "created_at,entities",
		})
		if err != nil {
			log.Fatalf("could not get tweets: %s", err.Error())
		}
		randomIndex := rand.Intn(len(tweets.Data))
		tweet := tweets.Data[randomIndex]
		tweetURL, err := getTweetURL(&tweet, user)
		if err != nil {
			log.Fatalf("failed to get tweet url %s", err.Error())
		}

		_, err = s.ChannelMessageSend(m.ChannelID, tweetURL.String())
		if err != nil {
			log.Fatalf("failed sending message to %s with content %s", m.ChannelID, tweetURL.String())
		}
	}
}

func getTweetURL(tweet *Tweet, user string) (*url.URL, error) {
	rawURL := fmt.Sprintf("https://twitter.com/%s/status/%s", user, tweet.ID)
	fmt.Println(user)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse url %s", err.Error())
	}

	return parsedURL, nil
}
