package utils

import (
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func StreamResponseChunks(api *slack.Client, channel, fullText string) {
	// Split into sentences or short chunks
	chunks := strings.Split(fullText, ". ")
	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		// Add a period back if missing
		if !strings.HasSuffix(chunk, ".") {
			chunk += "."
		}

		// Send message chunk
		api.PostMessage(channel, slack.MsgOptionText(chunk, false))

		// Simulate typing delay
		time.Sleep(1 * time.Second)
	}
}
