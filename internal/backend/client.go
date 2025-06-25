package backend

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// type ChatRequest struct {
// 	UserID string `json:"user_id"`
// 	Query  string `json:"query"`
// }

type ChatResponse struct {
	FullResponse string `json:"full_response"`
}

func SendQueryToBackend(url, userID, query string) (string, error) {
	reqBody := ChatRequest{
		UserID: userID,
		Query:  query,
	}

	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	var chatResp ChatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", err
	}

	return chatResp.FullResponse, nil
}
