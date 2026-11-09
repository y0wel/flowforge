package steps

import (
	"context"
	"errors"
	"strings"
	"testing"

	"flowforge/internal/llm"
)

func TestLLMStep_Execute_Success(t *testing.T) {
	fake := &llm.FakeLLM{Response: "invoice"}
	step := &LLMStep{Name: "classify", Prompt: "Classify this document", LLM: fake}
	data := NewContext()
	data.Data["input"] = "Invoice from Acme Corp"

	err := step.Execute(context.Background(), data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := data.Data["classify"].(string)
	if !ok {
		t.Fatalf("expected data.Data[%q] to be a string", step.Name)
	}
	if got != "invoice" {
		t.Errorf("expected %q, got %q", "invoice", got)
	}
}

func TestLLMStep_Execute_NoInput(t *testing.T) {
	fake := &llm.FakeLLM{Response: "ok"}
	step := &LLMStep{Name: "classify", Prompt: "Classify this document", LLM: fake}
	data := NewContext()

	if err := step.Execute(context.Background(), data); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLLMStep_Execute_LLMError(t *testing.T) {
	fake := &llm.FakeLLM{Err: errors.New("connection refused")}
	step := &LLMStep{Name: "classify", Prompt: "Classify this document", LLM: fake}
	data := NewContext()

	err := step.Execute(context.Background(), data)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "classify") {
		t.Errorf("expected error to mention step name, got %q", err.Error())
	}
}
