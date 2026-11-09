package steps

type Context struct {
	Data map[string]any
}

func NewContext() *Context {
	return &Context{
		Data: make(map[string]any),
	}
}
