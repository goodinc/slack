package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendWebhookSuccessful(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(sendWebhookHandlerOK))
	defer server.Close()

	testClient := NewWebhookClient(server.URL, 250*time.Millisecond)

	err := testClient.SendMessage(&Payload{Text: testText})
	if err != nil {
		t.Error(err)
	}
}

func sendWebhookHandlerOK(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("ok"))
}

func TestSendWebhookdMissingArgument(t *testing.T) {
	t.Parallel()

	testClient := NewWebhookClient("http://www.example.com", 250*time.Millisecond)

	err := testClient.SendMessage(nil)
	if err == nil {
		t.Error("should have failed, payload missing")
	}
}

func TestSendWebhookBadResponse(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(sendWebhookHandlerBad))
	defer server.Close()
	slackAddr = server.URL

	testClient := NewWebhookClient(server.URL, 250*time.Millisecond)

	err := testClient.SendMessage(&Payload{Text: testText})
	if err == nil {
		t.Error("should have failed, invalid channel")
	}
	if err.Error() != "some_error" {
		t.Error("received unexpected error")
	}
}

func sendWebhookHandlerBad(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("some_error"))
}

func TestWebhookMessageTimeout(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(sendMessageHandlerWait))
	defer server.Close()
	slackAddr = server.URL

	testClient := NewWebhookClient(server.URL, 250*time.Millisecond)

	err := testClient.SendMessage(&Payload{Text: testText})
	if err == nil {
		t.Error("should have failed, timeout")
	}
}

func TestIt(t *testing.T) {
	t.Parallel()
	testClient := NewWebhookClient("https://hooks.slack.com/services/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", DefaultTimeout)

	err := testClient.SendMessage(&Payload{Text: "This is a test"})
	if err != nil {
		t.Error(err)
	}
}
