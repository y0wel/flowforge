package steps

import (
	"fmt"

	"flowforge/internal/llm"
	"flowforge/internal/workflow"
)

func BuildStep(config workflow.StepConfig, provider llm.LLM) (Step, error) {
	switch config.Type {
	case "log":
		return &LogStep{
			Name:    config.Name,
			Message: config.Message,
		}, nil

	case "extract_text":
		return &ExtractTextStep{
			Name: config.Name,
		}, nil

	case "json_file":
		return &JSONFileStep{
			Name: config.Name,
			Path: config.Path,
		}, nil

	case "llm":
		return &LLMStep{
			Name:   config.Name,
			Prompt: config.Prompt,
			LLM:    provider,
		}, nil

	case "extract_json":
		return &ExtractJSONStep{
			Name:   config.Name,
			Prompt: config.Prompt,
			LLM:    provider,
		}, nil

	case "validate":
		return &ValidateStep{
			Name:     config.Name,
			Required: config.Required,
			Numeric:  config.Numeric,
		}, nil

	default:
		return nil, fmt.Errorf("unknown step type: %s", config.Type)
	}
}

func BuildSteps(wf workflow.Workflow, provider llm.LLM) ([]ConfiguredStep, error) {
	steps := make([]ConfiguredStep, 0, len(wf.Steps))

	for _, config := range wf.Steps {
		step, err := BuildStep(config, provider)
		if err != nil {
			return nil, fmt.Errorf("failed to build step %q: %w", config.Name, err)
		}
		steps = append(steps, ConfiguredStep{
			Name: config.Name,
			When: config.When,
			Step: step,
		})
	}
	return steps, nil
}
