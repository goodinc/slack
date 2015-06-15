package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Payload struct {
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconUrl     string       `json:"icon_url,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments",omitempty"`
}

type WebhookClient struct {
	webhook_url string
	timeout     time.Duration
}

const DefaultTimeout = time.Duration(10 * time.Second)

// NewWebhookClient returns a Client with the provided webhook url (default timeout to 10 seconds)
func NewWebhookClient(webhook string, timeout time.Duration) *WebhookClient {
	return &WebhookClient{webhook, timeout}
}

// SendMessage sends a text message to the default channel unless overridden
// https://api.slack.com/incoming-webhooks
func (c *WebhookClient) SendMessage(p *Payload) error {
	if p == nil {
		return errors.New("payload_missing")
	}

	client := http.Client{
		Timeout: time.Duration(c.timeout),
	}

	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	res, err := client.Post(c.webhook_url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	s := buf.String() // Does a complete copy of the bytes in the buffer.

	if s != "ok" {
		return errors.New(s)
	}

	return nil
}
