package workflow

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseWorkflow(filename string) (Workflow, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Workflow{}, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow Workflow
	err = yaml.Unmarshal(data, &workflow)
	if err != nil {
		return Workflow{}, fmt.Errorf("failed to parse workflow YAML: %w", err)
	}

	return workflow, nil
}
