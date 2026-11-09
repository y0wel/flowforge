package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"flowforge/internal/llm"
)

type ExtractJSONStep struct {
	Name   string
	Prompt string
	LLM    llm.LLM
}

func (s *ExtractJSONStep) Execute(ctx context.Context, data *Context) error {
	input, _ := data.Data["input"].(string)

	fullPrompt := s.Prompt
	if input != "" {
		fullPrompt = s.Prompt + "\n\n" + input
	}
	fullPrompt += "\n\nRespond with valid JSON only. Do not include any explanation or markdown formatting."

	response, err := s.LLM.Generate(ctx, fullPrompt)
	if err != nil {
		return fmt.Errorf("extract_json step %q: failed to generate response: %w", s.Name, err)
	}

	cleaned := stripJSONFences(response)

	var extracted map[string]any
	if err := json.Unmarshal([]byte(cleaned), &extracted); err != nil {
		return fmt.Errorf("extract_json step %q: invalid JSON response: %w (raw response: %s)", s.Name, err, response)
	}

	for key, value := range extracted {
		data.Data[key] = value
	}

	return nil
}

func stripJSONFences(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
