package workflow

type Workflow struct {
	Name  string       `yaml:"name"`
	Steps []StepConfig `yaml:"steps"`
}

type StepConfig struct {
	Name     string            `yaml:"name"`
	Type     string            `yaml:"type"`
	Prompt   string            `yaml:"prompt"`
	Message  string            `yaml:"message"`
	Path     string            `yaml:"path"`
	Required []string          `yaml:"required"`
	Numeric  []string          `yaml:"numeric"`
	When     map[string]string `yaml:"when"`
}
