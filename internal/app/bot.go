// package app

// import (
// 	"log"

// 	"github.com/manoj-2003/chatrelay/internal/config"
// 	slackhandler "github.com/manoj-2003/chatrelay/internal/slack"
// 	"github.com/slack-go/slack"
// 	"github.com/slack-go/slack/socketmode"
// )

// func RunBot(env *config.EnvVars) {
// 	api := slack.New(
// 		env.SlackBotToken,
// 		slack.OptionDebug(true),
// 		slack.OptionAppLevelToken(env.SlackAppToken),
// 	)

// 	client := socketmode.New(api, socketmode.OptionDebug(true))

// 	go slackhandler.HandleSlackEvents(api, client)

// 	log.Println("✅ Slack bot running via Socket Mode...")
// 	if err := client.Run(); err != nil {
// 		log.Fatalf("Socketmode client failed: %v", err)
// 	}
// }

package app

import (
	"log"

	"github.com/manoj-2003/chatrelay/internal/config"
	slackhandler "github.com/manoj-2003/chatrelay/internal/slack"
	telemetry "github.com/manoj-2003/chatrelay/internal/telementary"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func RunBot(env *config.EnvVars) {
	// 👇 Initialize OpenTelemetry Tracer
	shutdown := telemetry.InitTracer()
	defer shutdown()

	// ✅ Setup Slack client
	api := slack.New(
		env.SlackBotToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(env.SlackAppToken),
	)

	client := socketmode.New(api, socketmode.OptionDebug(true))

	// ✅ Handle events in goroutine
	go slackhandler.HandleSlackEvents(api, client)

	log.Println("✅ Slack bot running via Socket Mode...")
	if err := client.Run(); err != nil {
		log.Fatalf("Socketmode client failed: %v", err)
	}
}
