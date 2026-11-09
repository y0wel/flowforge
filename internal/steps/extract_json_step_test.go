package steps

import (
	"context"
	"errors"
	"testing"

	"flowforge/internal/llm"
)

func TestExtractJSONStep_Execute_Success(t *testing.T) {
	fake := &llm.FakeLLM{Response: `{"vendor": "Acme Corp", "total": 42.5}`}
	step := &ExtractJSONStep{Name: "extract", Prompt: "Extract fields", LLM: fake}
	data := NewContext()

	if err := step.Execute(context.Background(), data); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.Data["vendor"] != "Acme Corp" {
		t.Errorf("expected vendor %q, got %v", "Acme Corp", data.Data["vendor"])
	}
	if data.Data["total"] != 42.5 {
		t.Errorf("expected total 42.5, got %v", data.Data["total"])
	}
}

func TestExtractJSONStep_Execute_MarkdownFences(t *testing.T) {
	fake := &llm.FakeLLM{Response: "```json\n{\"vendor\": \"Acme Corp\"}\n```"}
	step := &ExtractJSONStep{Name: "extract", Prompt: "Extract fields", LLM: fake}
	data := NewContext()

	if err := step.Execute(context.Background(), data); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.Data["vendor"] != "Acme Corp" {
		t.Errorf("expected vendor %q, got %v", "Acme Corp", data.Data["vendor"])
	}
}

func TestExtractJSONStep_Execute_InvalidJSON(t *testing.T) {
	fake := &llm.FakeLLM{Response: "this is not json"}
	step := &ExtractJSONStep{Name: "extract", Prompt: "Extract fields", LLM: fake}
	data := NewContext()

	if err := step.Execute(context.Background(), data); err == nil {
		t.Fatal("expected an error for invalid JSON, got nil")
	}
}

func TestExtractJSONStep_Execute_LLMError(t *testing.T) {
	fake := &llm.FakeLLM{Err: errors.New("network error")}
	step := &ExtractJSONStep{Name: "extract", Prompt: "Extract fields", LLM: fake}
	data := NewContext()

	if err := step.Execute(context.Background(), data); err == nil {
		t.Fatal("expected an error, got nil")
	}
}
