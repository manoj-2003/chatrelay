package main

import (
	"fmt"

	"github.com/manoj-2003/chatrelay/internal/app"
	"github.com/manoj-2003/chatrelay/internal/config"
)

func main() {
	cfg := config.LoadEnv()

	fmt.Println("Slack App Token:", cfg.SlackAppToken)
	fmt.Println("Backend URL:", cfg.ChatBackendURL)
	fmt.Println("OpenAI Key:", cfg.OpenAIKey)
	app.RunBot(cfg)
}
