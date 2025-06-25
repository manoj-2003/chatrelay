package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	SlackAppToken   string
	SlackBotToken   string
	ChatBackendURL  string
	OpenAIKey       string
	OtelExporterURL string
	ServiceName     string
}

func LoadEnv() *EnvVars {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	return &EnvVars{
		SlackAppToken:   os.Getenv("SLACK_APP_TOKEN"),
		SlackBotToken:   os.Getenv("SLACK_BOT_TOKEN"),
		ChatBackendURL:  os.Getenv("CHAT_BACKEND_URL"),
		OpenAIKey:       os.Getenv("OPENAI_API_KEY"),
		OtelExporterURL: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		ServiceName:     os.Getenv("OTEL_SERVICE_NAME"),
	}
}
