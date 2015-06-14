// Copyright 2015 Bowery, Inc.

package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	slackAddr = "https://slack.com/api"
)

const (
	postMessageURI = "chat.postMessage"
)

type slackPostMessageRes struct {
	Ok    bool
	Error string
}

type Attachment struct {
	Fallback   string            `json:"fallback"`
	Color      string            `json:"color,omitempty"`
	Pretext    string            `json:"pretext,omitempty"`
	AuthorName string            `json:"author_name,omitempty"`
	AuthorLink string            `json:"author_link,omitempty"`
	AuthorIcon string            `json:"author_icon,omitempty"`
	Title      string            `json:"title,omitempty"`
	TitleLink  string            `json:"title_link,omitempty"`
	Text       string            `json:"text,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
	ImageURL   string            `json:"image_url,omitempty"`
	ThumbURL   string            `json:"thumb_url,omitempty"`
}

type AttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

// Client represents a slack api client. A Client is
// used for making requests to the slack api.
type Client struct {
	token    string
	username string
	icon     string

	timeout time.Duration
}

// NewClient returns a Client with the provided api token
// default timeout to 10 seconds
func NewClient(token, username, icon string) *Client {
	return &Client{token, username, icon, time.Duration(10 * time.Second)}
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SendMessage sends a text message to a specific channel
// with a specific username.
func (c *Client) SendMessage(channel, message string, attachments []Attachment) error {
	if channel == "" {
		return errors.New("channel required")
	}

	payload := url.Values{}
	payload.Set("token", c.token)
	payload.Set("username", c.username)
	payload.Set("icon_emoji", c.icon)
	payload.Set("channel", channel)
	payload.Set("text", message)

	if attachments != nil {
		d, err := json.Marshal(attachments)
		if err != nil {
			return err
		}

		payload.Set("attachments", string(d))
	}

	client := http.Client{
		Timeout: time.Duration(c.timeout),
	}

	res, err := client.PostForm(fmt.Sprintf("%s/%s", slackAddr, postMessageURI), payload)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody := new(slackPostMessageRes)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(resBody)
	if err != nil {
		return err
	}

	if !resBody.Ok {
		return errors.New(resBody.Error)
	}

	return nil
}
