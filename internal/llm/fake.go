package llm

import "context"

type FakeLLM struct {
	Response string
	Err      error
}

func (f *FakeLLM) Generate(ctx context.Context, prompt string) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.Response, nil
}
