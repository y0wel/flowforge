package steps

import (
	"context"
	"fmt"
)

type LogStep struct {
	Name    string
	Message string
}

func (step *LogStep) Execute(
	ctx context.Context,
	data *Context,
) error {
	fmt.Printf("[%s] %s\n", step.Name, step.Message)
	return nil
}
