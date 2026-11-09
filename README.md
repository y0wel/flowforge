# FlowForge

A small, extensible AI workflow engine written in Go. FlowForge runs workflows defined in YAML: a sequence of steps that can read input, call a local LLM, extract and validate structured data, branch conditionally, and write output.

> **This is a personal learning project.** It was built to learn Go (interfaces, composition, error handling, HTTP, JSON, and testing) by building a real (if small) application rather than isolated exercises. It is not a production-grade workflow engine, and intentionally keeps the architecture simple rather than "enterprise-ready."

---

## What it does

A workflow is a YAML file describing a named sequence of steps. Each step reads from and writes to a shared in-memory `Context` (a `map[string]any`), so later steps can use data produced by earlier ones. Steps can also be skipped conditionally based on values already in the context.

Core abstractions:

```go
type Step interface {
	Execute(ctx context.Context, data *Context) error
}

type LLM interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
```

Any step type just needs to implement `Step`. Any LLM backend just needs to implement `LLM` which is what lets FlowForge run fully offline in tests (`FakeLLM`) and against a real local model in practice, without the engine itself knowing the difference.

## Building and running

Requires Go 1.27+.

```bash
make build                         # builds ./flowforge
make run workflow=workflow.yml     # go run ./cmd/flowforge run workflow.yml
make test                          # go test ./...
make vet                           # go vet ./...
```

Or without `make`:

```bash
go build -o flowforge ./cmd/flowforge
./flowforge run workflow.yml
```

### CLI usage

```bash
flowforge run <workflow.yml>                      # no input
flowforge run <workflow.yml> <file>                # file contents become the workflow input
flowforge run <workflow.yml> --text "some text"    # inline text instead of a file
```

File and `--text` input are mutually exclusive. Either way, the input ends up in the workflow context under the `input` key, available to any step that needs it.

## Workflow YAML reference

```yaml
name: invoice_pipeline
steps:
  - name: extract
    type: extract_json
    prompt: |
      Extract information from this document and respond with exactly this JSON structure:
      { "vendor": "string", "invoice_number": "string", "total": 0.0, "currency": "string" }

  - name: validate
    type: validate
    required: [vendor, invoice_number, total, currency]
    numeric: [total]

  - name: save
    type: json_file
    path: ./output/invoice.json
    when:
      vendor: "Acme Corp"
```

Supported step types:

| Type | Config fields | Behavior |
| --- | --- | --- |
| `log` | `message` | Prints `[name] message`. |
| `extract_text` | — | Copies the workflow input into `text`. |
| `llm` | `prompt` | Sends `prompt` (plus the input, if any) to the configured LLM; stores the raw text response under the step's `name`. |
| `extract_json` | `prompt` | Like `llm`, but asks for JSON and parses the response, merging the resulting fields directly into the context. |
| `validate` | `required`, `numeric` | Fails the workflow if any `required` field is missing/empty, or any `numeric` field isn't a number. |
| `json_file` | `path` | Serializes the entire context to indented JSON and writes it to `path` (creating directories as needed). |

Any step may include a `when:` block (`field: expected-value`) — the step only runs if every listed context field matches exactly; otherwise it's skipped.

## Configuring an LLM provider

FlowForge talks to any OpenAI-compatible chat-completions server (`POST /v1/chat/completions`). This has been tested against [LM Studio](https://lmstudio.ai/) and Apple's on-device model via [apfel](https://apfel.franzai.com/). The provider is wired up in [cmd/flowforge/main.go](cmd/flowforge/main.go):

```go
provider := llm.NewLMStudioLLM("http://localhost:11434", "apple-foundationmodel")
```

Point `baseURL`/`model` at whichever local server and model you're running. Tests never hit a real server. They use `llm.FakeLLM`, which returns a canned response or error.

## Project structure

```text
cmd/flowforge/         main.go — CLI entry point, wiring only
internal/workflow/     YAML schema (Workflow, StepConfig), parsing, workflow-level validation
internal/llm/          LLM interface, FakeLLM, LMStudioLLM (OpenAI-compatible client)
internal/steps/        Step interface, Context, Executor, and all concrete step implementations
```

## Status

The core engine, CLI, all six step types, conditional execution, an LLM integration, and unit tests are implemented and working end-to-end. Remaining ideas (structured logging, execution IDs, batch/concurrent processing, an HTTP API, etc.) are tracked as GitHub issues rather than in this file.
