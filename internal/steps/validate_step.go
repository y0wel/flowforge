package steps

import (
	"context"
	"fmt"
)

type ValidateStep struct {
	Name     string
	Required []string
	Numeric  []string
}

func (s *ValidateStep) Execute(ctx context.Context, data *Context) error {
	for _, field := range s.Required {
		value, ok := data.Data[field]
		if !ok {
			return fmt.Errorf("validate step %q: missing required field %q", s.Name, field)
		}
		if str, isStr := value.(string); isStr && str == "" {
			return fmt.Errorf("validate step %q: field %q is empty", s.Name, field)
		}
	}

	for _, field := range s.Numeric {
		value, ok := data.Data[field]
		if !ok {
			continue
		}
		if _, isNum := value.(float64); !isNum {
			return fmt.Errorf("validate step %q: field %q must be a number, got %T", s.Name, field, value)
		}
	}

	return nil
}
