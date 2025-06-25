package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ChatRequest struct {
	UserID string `json:"user_id"`
	Query  string `json:"query"`
}

type ChatResponse struct {
	FullResponse string `json:"full_response"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("🧠 Received query from user %s: %s", req.UserID, req.Query)

	response := ChatResponse{
		FullResponse: "Goroutines are lightweight, concurrent execution units in Go. They allow for massive parallelism.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/v1/chat/stream", handler)
	log.Println("🚀 Mock backend running at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
