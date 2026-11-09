package steps

type ConfiguredStep struct {
	Name string
	When map[string]string
	Step Step
}

func matchesCondition(when map[string]string, data *Context) bool {
	for key, expected := range when {
		actual, ok := data.Data[key].(string)
		if !ok || actual != expected {
			return false
		}
	}
	return true
}
