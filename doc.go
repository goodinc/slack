// Package slack is a client library for Slack.
//
// Usage:
/*
  package main

  import (
    "github.com/Bowery/slack"
  )

  var (
    client *slack.Client
  )

  func main() {
    client = slack.NewClient("API_TOKEN")
    err := client.SendMessage("#mychannel", "message", "username")
    if err != nil {
      log.Fatal(err)
    }

    whClient := NewWebhookClient("https://hooks.slack.com/services/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", DefaultTimeout)
    err := whClient.SendMessage(&Payload{Text: "This is a test"})
    if err != nil {
      log.Fatal(err)
    }
  }
*/
package slack
