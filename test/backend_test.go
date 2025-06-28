package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/manoj-2003/chatrelay/internal/backend"
)

func TestStreamQueryToGroq_Mock(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected http.ResponseWriter to be Flusher")
		}

		// Simulate stream data
		fmt.Fprint(w, "data: Hello\n\n")
		flusher.Flush()
		fmt.Fprint(w, "data: World\n\n")
		flusher.Flush()
	}))
	defer mock.Close()

	streamChan := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go backend.StreamQueryToGroq(ctx, mock.URL, "Hello", "test", func(msg string) {
		streamChan <- msg
	})

	var result string
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case msg, ok := <-streamChan:
			if !ok {
				break loop
			}
			result += msg
		}
	}

	if !strings.Contains(result, "Hello") || !strings.Contains(result, "World") {
		t.Errorf("Unexpected stream result: %q", result)
	}
}
