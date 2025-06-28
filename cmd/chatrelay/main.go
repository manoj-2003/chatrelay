package main

import (
	"context"
	"fmt"
	"time"

	"github.com/manoj-2003/chatrelay/internal/app"
	"go.opentelemetry.io/otel"

	"github.com/manoj-2003/chatrelay/internal/config"
	telemetry "github.com/manoj-2003/chatrelay/internal/telementary"
)

func main() {
	cfg := config.LoadEnv()
	shutdown := telemetry.InitTracer()
	defer shutdown()

	fmt.Println("Slack App Token:", cfg.SlackAppToken)
	fmt.Println("Backend URL:", cfg.ChatBackendURL)
	fmt.Println("GROQ API Key:", cfg.Grok)

	// ✅ Force a test span and flush before bot runs
	tr := otel.Tracer("chatrelay")
	_, span := tr.Start(context.Background(), "bot-start-span")
	span.End()

	// ⚠️ Add sleep to give exporter time to send the span before app shuts down
	time.Sleep(3 * time.Second) // 2–3s is usually enough

	app.RunBot(cfg)
}
