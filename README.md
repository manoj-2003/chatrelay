# chatrelay
# ChatRelay Bot with OpenTelemetry and Jaeger

This project is a Slack bot (`chatrelay`) instrumented with OpenTelemetry for distributed tracing, using Docker Compose to run the bot, OpenTelemetry Collector, and Jaeger for trace visualization.

---

## Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/)
---

## Quick Start

1. **Clone the repository and navigate to the project directory:**
   ```sh
   git clone <your-repo-url>
   cd chatrelay
   ```

2. **Configure environment variables:**
   - Edit the `.env` file with your Slack tokens and API keys as needed.

3. **Start all services:**
   ```sh
   docker-compose up --build
   ```

4. **Access Jaeger UI:**
   - Open [http://localhost:16686](http://localhost:16686) in your browser to view traces.

5. **Use the bot:**
   - Interact with your Slack bot as usual. Traces will appear in Jaeger under the `chatrelay` service.

---

## Environment Variables

Set in `.env`:

- `SLACK_APP_TOKEN` – Your Slack app token
- `SLACK_BOT_TOKEN` – Your Slack bot token
- `CHAT_BACKEND_URL` – Backend URL for chat relay
- `OTEL_EXPORTER_OTLP_ENDPOINT` – Should be `otel-collector:4317`
- `OTEL_SERVICE_NAME` – Service name for tracing (e.g., `chatrelay`)
- `GROQ_API_KEY` – Your GROQ API key

---

## Troubleshooting

- If you see connection errors to `otel-collector`, ensure all services are running via Docker Compose and the endpoint is set to `otel-collector:4317`.
- For tracing to work, your Go code must use the OTLP gRPC exporter (`otlptracegrpc`) and the collector must listen on `0.0.0.0:4317`.

---
