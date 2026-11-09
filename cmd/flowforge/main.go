package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"flowforge/internal/llm"
	"flowforge/internal/steps"
	"flowforge/internal/workflow"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: missing command")
		fmt.Fprintln(os.Stderr, "Usage: flowforge run <workflow.yml>")
		os.Exit(1)
	}

	command := args[1]

	if command != "run" {
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n", command)
		fmt.Fprintln(os.Stderr, "Usage: flowforge run <workflow.yml>")
		os.Exit(1)
	}

	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: missing workflow file")
		fmt.Fprintln(os.Stderr, "Usage: flowforge run <workflow.yml>")
		os.Exit(1)
	}

	if len(args) > 5 {
		fmt.Fprintln(os.Stderr, "Error: too many arguments")
		fmt.Fprintln(os.Stderr, "Usage: flowforge run <workflow.yml>")
		os.Exit(1)
	}

	inputFile := ""
	inputText := ""
	contents := []byte{}

	if len(args) >= 4 {
		if args[3] == "--text" {
			if len(args) != 5 {
				fmt.Fprintln(os.Stderr, "Error: --text requires a value")
				fmt.Fprintln(os.Stderr, "Usage: flowforge run <workflow.yml> --text \"...\"")
				os.Exit(1)
			}
			inputText = args[4]
		} else {
			if len(args) != 4 {
				fmt.Fprintln(os.Stderr, "Error: too many arguments")
				fmt.Fprintln(os.Stderr, "Usage: flowforge run <workflow.yml> [file]")
				os.Exit(1)
			}
			inputFile = args[3]

			file, err := os.Open(inputFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error opening input file:", err)
				os.Exit(1)
			}
			defer file.Close()

			contents, err = io.ReadAll(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input file:", err)
				os.Exit(1)
			}
		}
	}

	input := inputText
	if input == "" {
		input = string(contents)
	}

	workflowFile := args[2]

	wf, err := workflow.ParseWorkflow(workflowFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if err := workflow.ValidateWorkflow(wf); err != nil {
		fmt.Fprintln(os.Stderr, "Invalid workflow:")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// provider := &llm.FakeLLM{Response: "fake response"}
	// Use llm.NewLMStudioLLM as the LLM implementation
	// provider := llm.NewLMStudioLLM("http://localhost:1234", "my-model")
	provider := llm.NewLMStudioLLM("http://localhost:11434", "apple-foundationmodel")

	builtSteps, err := steps.BuildSteps(wf, provider)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Running %s\n", workflowFile)
	fmt.Printf("Workflow: %s\n", wf.Name)
	// Here you would add the logic to run the workflow defined in the YAML file.

	data := steps.NewContext()
	data.Data["input"] = input

	executor := steps.NewExecutor(builtSteps)

	err = executor.Run(context.Background(), data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running workflow:", err)
		os.Exit(1)
	}

	fmt.Println("Workflow executed successfully.")
	fmt.Println("Extracted Data:", data.Data)

	for _, step := range wf.Steps {
		fmt.Printf("Step Name: %s, Type: %s\n", step.Name, step.Type)
	}
}
