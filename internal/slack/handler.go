// package slackhandler

// import (
// 	"context"
// 	"log"
// 	"strings"

// 	"github.com/manoj-2003/chatrelay/internal/utils"

// 	"github.com/slack-go/slack"
// 	"github.com/slack-go/slack/slackevents"
// 	"github.com/slack-go/slack/socketmode"

// 	"github.com/manoj-2003/chatrelay/internal/backend"
// 	"github.com/manoj-2003/chatrelay/internal/config"
// )

// func HandleSlackEvents(api *slack.Client, client *socketmode.Client) {
// 	env := config.LoadEnv()

// 	for evt := range client.Events {
// 		switch evt.Type {
// 		case socketmode.EventTypeEventsAPI:
// 			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
// 			if !ok {
// 				log.Printf("⚠️ Could not cast to EventsAPIEvent: %+v\n", evt)
// 				continue
// 			}

// 			client.Ack(*evt.Request)

// 			switch event := eventsAPIEvent.InnerEvent.Data.(type) {

// 			case *slackevents.AppMentionEvent:
// 				log.Printf("👋 Mention from user %s: %s", event.User, event.Text)

// 				// ✅ Correct argument order

// 				ctx := context.Background()
// 				reply, err := backend.SendQueryToGroq(ctx, event.User, event.Text, env.Grok)
// 				if err != nil {
// 					log.Printf("❌ Groq error: %v", err) // <-- ADD THIS FOR DEBUGGING
// 					reply = ":x: " + err.Error()        // <-- SHOW actual Groq error in Slack
// 				}

// 				utils.StreamResponseChunks(api, event.Channel, reply)

// 			case *slackevents.MessageEvent:
// 				// ✅ Only reply if not a bot and it's a DM
// 				if event.BotID == "" && strings.HasPrefix(event.Channel, "D") {
// 					log.Printf("💬 DM from %s: %s", event.User, event.Text)

// 					ctx := context.Background()
// 					reply, err := backend.SendQueryToGroq(ctx, event.User, event.Text, env.Grok)

// 					if err != nil {
// 						log.Printf("❌ Groq error: %v", err) // <-- ADD THIS FOR DEBUGGING
// 						reply = ":x: " + err.Error()        // <-- SHOW actual Groq error in Slack
// 					}

// 					api.PostMessage(event.Channel, slack.MsgOptionText(reply, false))
// 				}
// 			}

// 		default:
// 			log.Printf("📦 Unhandled event type: %s\n", evt.Type)
// 		}
// 	}
// }

package slackhandler

import (
	"context"
	"log"
	"strings"

	"github.com/manoj-2003/chatrelay/internal/backend"
	"github.com/manoj-2003/chatrelay/internal/config"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("chatrelay/slack")

func HandleSlackEvents(api *slack.Client, client *socketmode.Client) {
	env := config.LoadEnv()

	for evt := range client.Events {
		switch evt.Type {
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				log.Printf("⚠️ Could not cast to EventsAPIEvent: %+v\n", evt)
				continue
			}

			client.Ack(*evt.Request)

			switch event := eventsAPIEvent.InnerEvent.Data.(type) {

			case *slackevents.AppMentionEvent:
				ctx, span := tracer.Start(context.Background(), "AppMentionEvent")
				span.SetAttributes(
					attribute.String("user_id", event.User),
					attribute.String("query", event.Text),
				)
				log.Printf("👋 Mention from user %s: %s", event.User, event.Text)

				err := backend.StreamQueryToGroq(ctx, event.User, event.Text, env.Grok, func(line string) {
					line = strings.TrimSpace(line)
					if line != "" {
						_, _, postErr := api.PostMessage(event.Channel, slack.MsgOptionText(line, false))
						if postErr != nil {
							log.Printf("❌ Failed to post streamed line: %v", postErr)
						}
					}
				})
				if err != nil {
					log.Printf("❌ Groq stream error (trace %s): %v", span.SpanContext().TraceID(), err)
					api.PostMessage(event.Channel, slack.MsgOptionText(":x: "+err.Error(), false))
				}

				span.End()

			case *slackevents.MessageEvent:
				if event.BotID == "" && strings.HasPrefix(event.Channel, "D") {
					ctx, span := tracer.Start(context.Background(), "DirectMessageEvent")
					span.SetAttributes(
						attribute.String("user_id", event.User),
						attribute.String("query", event.Text),
					)
					log.Printf("💬 DM from %s: %s", event.User, event.Text)

					// Optional: Indicate the bot is thinking
					api.PostMessage(event.Channel, slack.MsgOptionText("_Thinking..._", false))

					err := backend.StreamQueryToGroq(ctx, event.User, event.Text, env.Grok, func(line string) {
						line = strings.TrimSpace(line)
						if line != "" {
							_, _, postErr := api.PostMessage(event.Channel, slack.MsgOptionText(line, false))
							if postErr != nil {
								log.Printf("❌ Failed to post streamed line: %v", postErr)
							}
						}
					})
					if err != nil {
						log.Printf("❌ Groq stream error (trace %s): %v", span.SpanContext().TraceID(), err)
						api.PostMessage(event.Channel, slack.MsgOptionText(":x: "+err.Error(), false))
					}

					span.End()
				}
			}

		default:
			log.Printf("📦 Unhandled event type: %s\n", evt.Type)
		}
	}
}
