package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LMStudioLLM struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

func NewLMStudioLLM(baseURL, model string) *LMStudioLLM {
	return &LMStudioLLM{
		BaseURL: baseURL,
		Model:   model,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type lmStudioMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type lmStudioRequest struct {
	Model    string            `json:"model"`
	Messages []lmStudioMessage `json:"messages"`
}

type lmStudioResponse struct {
	Choices []struct {
		Message lmStudioMessage `json:"message"`
	} `json:"choices"`
}

func (l *LMStudioLLM) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody := lmStudioRequest{
		Model: l.Model,
		Messages: []lmStudioMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("lmstudio: failed to encode request: %w", err)
	}

	url := l.BaseURL + "/v1/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("lmstudio: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("lmstudio: request failed (is the server running?): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("lmstudio: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result lmStudioResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("lmstudio: failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("lmstudio: no choices returned")
	}

	return result.Choices[0].Message.Content, nil
}
