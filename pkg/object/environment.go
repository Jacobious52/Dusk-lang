package object

// Environment stores the variables for a context
type Environment struct {
	vars   map[string]Object
	parent *Environment
}

// NewEnvironment makes a new Environment with no values
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{s, nil}
}

// NewChildEnvironment creates an enclosed Environment on the parent
func NewChildEnvironment(parent *Environment) *Environment {
	env := NewEnvironment()
	env.parent = parent
	return env
}

// Get a value from the variables map
func (e *Environment) Get(name string) (Object, bool) {
	o, ok := e.vars[name]
	if !ok && e.parent != nil {
		o, ok = e.parent.Get(name)
	}
	return o, ok
}

// Set a value in the varibles map
func (e *Environment) Set(name string, val Object) Object {
	e.vars[name] = val
	return val
}

// Assign a existing value in the varibles map
func (e *Environment) Assign(name string, val Object) (Object, bool) {
	_, ok := e.vars[name]
	if !ok && e.parent != nil {
		return e.parent.Assign(name, val)
	}
	e.vars[name] = val
	return val, ok
}
