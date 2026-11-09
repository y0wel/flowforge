package steps

import (
	"context"
	"fmt"

	"flowforge/internal/llm"
)

type LLMStep struct {
	Name   string
	Prompt string
	LLM    llm.LLM
}

func (s *LLMStep) Execute(ctx context.Context, data *Context) error {
	input, _ := data.Data["input"].(string)

	fullPrompt := s.Prompt
	if input != "" {
		fullPrompt = s.Prompt + "\n\n" + input
	}

	response, err := s.LLM.Generate(ctx, fullPrompt)
	if err != nil {
		return fmt.Errorf("llm step %q: failed to generate response: %w", s.Name, err)
	}
	data.Data[s.Name] = response
	return nil
}
