// package utils

// import (
// 	"strings"
// 	"time"

// 	"github.com/slack-go/slack"
// )

// func StreamResponseChunks(api *slack.Client, channel, fullText string) {
// 	// Split into sentences or short chunks
// 	chunks := strings.Split(fullText, ". ")
// 	for _, chunk := range chunks {
// 		chunk = strings.TrimSpace(chunk)
// 		if chunk == "" {
// 			continue
// 		}

// 		// Add a period back if missing
// 		if !strings.HasSuffix(chunk, ".") {
// 			chunk += "."
// 		}

// 		// Send message chunk
// 		api.PostMessage(channel, slack.MsgOptionText(chunk, false))

// 		// Simulate typing delay
// 		time.Sleep(1 * time.Second)
// 	}
// }

package utils

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("chatrelay/utils")

func StreamResponseChunks(ctx context.Context, api *slack.Client, channel, fullText string) {
	ctx, span := tracer.Start(ctx, "StreamResponseChunks")
	defer span.End()

	const chunkSize = 4000 // Slack max message size is around 4000
	paragraphs := splitIntoChunks(fullText, chunkSize)

	for i, chunk := range paragraphs {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		time.Sleep(1 * time.Second) // Simulate typing delay

		_, _, err := api.PostMessage(channel, slack.MsgOptionText(chunk, false))
		if err != nil {
			log.Printf("❌ Failed to send chunk %d: %v", i, err)
			span.RecordError(err)
			continue
		}

		// Add event to trace
		span.AddEvent("Sent message chunk", trace.WithAttributes(
			attribute.Int("chunk.index", i),
			attribute.Int("chunk.length", len(chunk)),
		))
	}
}

func splitIntoChunks(text string, size int) []string {
	var chunks []string
	for len(text) > size {
		cut := size
		if idx := strings.LastIndex(text[:size], " "); idx > 0 {
			cut = idx
		}
		chunks = append(chunks, text[:cut])
		text = text[cut:]
	}
	if len(text) > 0 {
		chunks = append(chunks, text)
	}
	return chunks
}
