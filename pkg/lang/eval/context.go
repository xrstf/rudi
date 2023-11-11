package eval

type Context struct {
	Object    Object
	Variables map[string]interface{}
}

func NewContext() Context {
	return Context{
		Variables: map[string]interface{}{},
	}
}

func (c Context) WithVariable(name string, val interface{}) Context {
	newVariables := map[string]interface{}{}
	for k, v := range c.Variables {
		newVariables[k] = v
	}
	newVariables[name] = val

	return Context{
		Object:    c.Object, // deep copy it? yeah, right?
		Variables: newVariables,
	}
}
