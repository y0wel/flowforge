package steps

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLogStep_Execute(t *testing.T) {
	step := &LogStep{Name: "hello", Message: "Hello World"}
	data := NewContext()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := step.Execute(context.Background(), data)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Hello World") {
		t.Errorf("expected output to contain message, got %q", output)
	}
	if !strings.Contains(output, "[hello]") {
		t.Errorf("expected output to contain step name, got %q", output)
	}
}
