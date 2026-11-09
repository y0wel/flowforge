package workflow

import "fmt"

func ValidateWorkflow(workflow Workflow) error {
	if len(workflow.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	for _, step := range workflow.Steps {
		if step.Name == "" {
			return fmt.Errorf("step name cannot be empty")
		}

		if step.Type == "" {
			return fmt.Errorf("step %q has no type defined", step.Name)
		}

		if !isValidStepType(step.Type) {
			return fmt.Errorf("step %q has an invalid type: %q", step.Name, step.Type)
		}

		if step.Type == "log" && step.Message == "" {
			return fmt.Errorf("step %q:\n log step requires a message", step.Name)
		}

		if step.Type == "json_file" && step.Path == "" {
			return fmt.Errorf("step %q:\n json_file step requires a path", step.Name)
		}

		if step.Type == "llm" && step.Prompt == "" {
			return fmt.Errorf("step %q:\n llm step requires a prompt", step.Name)
		}

		if step.Type == "extract_json" && step.Prompt == "" {
			return fmt.Errorf("step %q:\n extract_json step requires a prompt", step.Name)
		}
	}
	return nil
}

func isValidStepType(stepType string) bool {
	switch stepType {
	case "extract_text", "log", "json_file", "llm", "extract_json", "validate":
		return true
	default:
		return false
	}
}
