package steps

import (
	"context"
	"fmt"
)

type Executor struct {
	steps []ConfiguredStep
}

func NewExecutor(steps []ConfiguredStep) *Executor {
	return &Executor{
		steps: steps,
	}
}

func (e *Executor) Run(ctx context.Context, data *Context) error {
	for _, step := range e.steps {
		if !matchesCondition(step.When, data) {
			fmt.Printf("[%s] skipped (condition not met)\n", step.Name)
			continue
		}
		if err := step.Step.Execute(ctx, data); err != nil {
			return fmt.Errorf("step execution failed: %w", err)
		}
	}
	return nil
}
