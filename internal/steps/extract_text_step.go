package steps

import (
	"context"
	"fmt"
)

type ExtractTextStep struct {
	Name string
}

func (step *ExtractTextStep) Execute(
	ctx context.Context,
	data *Context,
) error {
	input, ok := data.Data["input"].(string)
	if !ok || input == "" {
		return fmt.Errorf("extract_text step %q: no input provided", step.Name)
	}
	data.Data["text"] = input
	return nil
}
