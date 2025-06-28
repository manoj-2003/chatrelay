package test

import (
	"testing"

	"github.com/slack-go/slack/slackevents"
)

func TestHandleAppMention(t *testing.T) {
	event := &slackevents.AppMentionEvent{
		User:    "U12345",
		Text:    "<@U12345> what is AI?",
		Channel: "C12345",
	}
	t.Logf("Testing AppMentionEvent: %+v", event)
}
