package object

// Environment stores the variables for a context
type Environment struct {
	vars map[string]Object
}

// NewEnvironment makes a new Environment with no values
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{s}
}

// Get a value from the variables map
func (e *Environment) Get(name string) (Object, bool) {
	o, ok := e.vars[name]
	return o, ok
}

// Set a value in the varibles map
func (e *Environment) Set(name string, val Object) Object {
	e.vars[name] = val
	return val
}
