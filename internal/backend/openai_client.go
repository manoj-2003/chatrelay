package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ChatRequest struct {
	UserID string `json:"user_id"`
	Query  string `json:"query"`
}

func SendQueryToOpenAI(ctx context.Context, user, query, backendURL string) (string, error) {
	payload := ChatRequest{
		UserID: user,
		Query:  query,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", backendURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error contacting backend: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read response as needed (stream or full)
	// For now, just read the whole body as a string
	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	return respBody.String(), nil
}

// This function is not part of the original code, but is added here to show the usage of SendQueryToOpenAI
