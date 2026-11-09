package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type JSONFileStep struct {
	Name string
	Path string
}

func (s *JSONFileStep) Execute(ctx context.Context, data *Context) error {
	jsonBytes, err := json.MarshalIndent(data.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("json_file step %q: failed to serialize data: %w", s.Name, err)
	}

	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("json_file step %q: failed to create output directory: %w", s.Name, err)
	}

	if err := os.WriteFile(s.Path, jsonBytes, 0o644); err != nil {
		return fmt.Errorf("json_file step %q: failed to write file: %w", s.Name, err)
	}

	return nil
}
