package slackhandler

import (
	"context"
	"log"
	"strings"

	"github.com/manoj-2003/chatrelay/internal/utils"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/manoj-2003/chatrelay/internal/backend"
	"github.com/manoj-2003/chatrelay/internal/config"
)

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
				log.Printf("👋 Mention from user %s: %s", event.User, event.Text)

				// ✅ Correct argument order

				ctx := context.Background()
				reply, err := backend.SendQueryToOpenAI(ctx, event.User, event.Text, env.ChatBackendURL)

				if err != nil {
					reply = "❌ Failed to contact backend"
				}

				utils.StreamResponseChunks(api, event.Channel, reply)

			case *slackevents.MessageEvent:
				// ✅ Only reply if not a bot and it's a DM
				if event.BotID == "" && strings.HasPrefix(event.Channel, "D") {
					log.Printf("💬 DM from %s: %s", event.User, event.Text)

					ctx := context.Background()
					reply, err := backend.SendQueryToOpenAI(ctx, event.User, event.Text, env.ChatBackendURL)

					if err != nil {
						reply = "❌ Failed to contact backend"
					}

					api.PostMessage(event.Channel, slack.MsgOptionText(reply, false))
				}
			}

		default:
			log.Printf("📦 Unhandled event type: %s\n", evt.Type)
		}
	}
}
