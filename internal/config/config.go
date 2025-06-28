package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	SlackAppToken               string
	SlackBotToken               string
	ChatBackendURL              string
	Grok                        string
	OtelExporterURL             string
	ServiceName                 string
	OTEL_EXPORTER_OTLP_ENDPOINT string
}

func LoadEnv() *EnvVars {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	return &EnvVars{
		SlackAppToken:               os.Getenv("SLACK_APP_TOKEN"),
		SlackBotToken:               os.Getenv("SLACK_BOT_TOKEN"),
		ChatBackendURL:              os.Getenv("CHAT_BACKEND_URL"),
		Grok:                        os.Getenv("GROQ_API_KEY"),
		OtelExporterURL:             os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		ServiceName:                 os.Getenv("OTEL_SERVICE_NAME"),
		OTEL_EXPORTER_OTLP_ENDPOINT: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
	}
}
