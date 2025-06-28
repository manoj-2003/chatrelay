// package backend

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// )

// type GroqRequest struct {
// 	Model    string    `json:"model"`
// 	Messages []Message `json:"messages"`
// }

// type Message struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }

// type GroqResponse struct {
// 	Choices []struct {
// 		Message struct {
// 			Content string `json:"content"`
// 		} `json:"message"`
// 	} `json:"choices"`
// 	Error *GroqError `json:"error,omitempty"` // added to catch error responses
// }

// type GroqError struct {
// 	Message string `json:"message"`
// 	Type    string `json:"type"`
// 	Param   string `json:"param,omitempty"`
// 	Code    string `json:"code,omitempty"`
// }

// func SendQueryToGroq(ctx context.Context, userID, query, apiKey string) (string, error) {
// 	requestData := GroqRequest{
// 		Model: "llama3-70b-8192",
// 		Messages: []Message{
// 			{Role: "user", Content: query},
// 		},
// 	}

// 	payload, err := json.Marshal(requestData)
// 	if err != nil {
// 		return "", fmt.Errorf("❌ Failed to marshal request: %v", err)
// 	}

// 	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(payload))
// 	if err != nil {
// 		return "", fmt.Errorf("❌ Failed to create request: %v", err)
// 	}

// 	req.Header.Set("Authorization", "Bearer "+apiKey)
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("❌ Failed to call Groq API: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, _ := io.ReadAll(resp.Body)

// 	// Parse the response into a struct, whether success or error
// 	var result GroqResponse
// 	if err := json.Unmarshal(bodyBytes, &result); err != nil {
// 		return "", fmt.Errorf("❌ Failed to parse response: %v\nRaw: %s", err, bodyBytes)
// 	}

// 	// If there's an error object in the response
// 	if result.Error != nil {
// 		return "", fmt.Errorf("❌ Groq API Error: %s (%s)", result.Error.Message, result.Error.Type)
// 	}

// 	if len(result.Choices) == 0 {
// 		return "🤖 No response from Groq", nil
// 	}

// 	return result.Choices[0].Message.Content, nil
// }

// package backend

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// )

// type GroqRequest struct {
// 	Model    string    `json:"model"`
// 	Messages []Message `json:"messages"`
// }

// type Message struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }

// type GroqResponse struct {
// 	Choices []struct {
// 		Message struct {
// 			Content string `json:"content"`
// 		} `json:"message"`
// 	} `json:"choices"`
// 	Error *GroqError `json:"error,omitempty"`
// }

// type GroqError struct {
// 	Message string `json:"message"`
// 	Type    string `json:"type"`
// 	Param   string `json:"param,omitempty"`
// 	Code    string `json:"code,omitempty"`
// }

// var tracer = otel.Tracer("chatrelay/backend")

// func SendQueryToGroq(ctx context.Context, userID, query, apiKey string) (string, error) {
// 	ctx, span := tracer.Start(ctx, "SendQueryToGroq")
// 	defer span.End()

// 	span.SetAttributes(
// 		attribute.String("user.id", userID),
// 		attribute.String("query", query),
// 	)

// 	requestData := GroqRequest{
// 		Model: "llama3-70b-8192",
// 		Messages: []Message{
// 			{Role: "user", Content: query},
// 		},
// 	}

// 	payload, err := json.Marshal(requestData)
// 	if err != nil {
// 		span.RecordError(err)
// 		return "", fmt.Errorf("❌ Failed to marshal request: %v", err)
// 	}

// 	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(payload))
// 	if err != nil {
// 		span.RecordError(err)
// 		return "", fmt.Errorf("❌ Failed to create request: %v", err)
// 	}

// 	req.Header.Set("Authorization", "Bearer "+apiKey)
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		span.RecordError(err)
// 		return "", fmt.Errorf("❌ Failed to call Groq API: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, _ := io.ReadAll(resp.Body)

// 	var result GroqResponse
// 	if err := json.Unmarshal(bodyBytes, &result); err != nil {
// 		span.RecordError(err)
// 		return "", fmt.Errorf("❌ Failed to parse response: %v\nRaw: %s", err, bodyBytes)
// 	}

// 	if result.Error != nil {
// 		err := fmt.Errorf("❌ Groq API Error: %s (%s)", result.Error.Message, result.Error.Type)
// 		span.RecordError(err)
// 		return "", err
// 	}

// 	if len(result.Choices) == 0 {
// 		return "🤖 No response from Groq", nil
// 	}

// 	reply := result.Choices[0].Message.Content
// 	span.SetAttributes(attribute.String("response", reply))

// 	return reply, nil
// }

package backend

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type GroqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"` // ✅ Enable streaming
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

var tracer = otel.Tracer("chatrelay/backend")

// 🚀 Stream each response line to a callback
func StreamQueryToGroq(ctx context.Context, userID, query, apiKey string, onLine func(string)) error {
	ctx, span := tracer.Start(ctx, "StreamQueryToGroq")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.String("query", query),
	)

	// ✅ Enable streaming
	requestData := GroqRequest{
		Model: "llama3-70b-8192",
		Messages: []Message{
			{Role: "user", Content: query},
		},
		Stream: true,
	}

	payload, err := json.Marshal(requestData)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("❌ Failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("❌ Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("❌ Failed to call Groq API: %v", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var buffer strings.Builder

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			span.RecordError(err)
			return fmt.Errorf("❌ Failed to read response: %v", err)
		}

		// Only process "data: " lines
		if !bytes.HasPrefix(line, []byte("data: ")) {
			continue
		}

		data := bytes.TrimPrefix(line, []byte("data: "))
		data = bytes.TrimSpace(data)

		if string(data) == "[DONE]" {
			break
		}

		var chunk GroqStreamChunk
		if err := json.Unmarshal(data, &chunk); err != nil {
			span.RecordError(err)
			continue
		}

		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				buffer.WriteString(content)

				// Stream full lines only
				if strings.Contains(content, "\n") {
					lines := strings.Split(buffer.String(), "\n")
					for i := 0; i < len(lines)-1; i++ {
						onLine(strings.TrimSpace(lines[i]))
					}
					buffer.Reset()
					buffer.WriteString(lines[len(lines)-1])
				}
			}
		}
	}

	// Flush any remaining content
	remaining := strings.TrimSpace(buffer.String())
	if remaining != "" {
		onLine(remaining)
	}

	return nil
}
