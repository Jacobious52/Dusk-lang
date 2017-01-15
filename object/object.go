package object

import "fmt"

// Type is a flag for the type of object
type Type int

const (
	NilType     Type = iota // NilType nil
	IntType                 // IntType int64
	FloatType               // FloatType float64
	BooleanType             // BooleanType bool
	StringType              // StringType string
)

// String for type
func (t Type) String() string {
	switch t {
	case IntType:
		return "int"
	case FloatType:
		return "float"
	case BooleanType:
		return "bool"
	case StringType:
		return "string"
	default:
		return "unknown"
	}
}

// Object is the base interface of all objects
type Object interface {
	Type() Type
	String() string
}

// Integer is a int64
type Integer struct {
	Value int64
}

// String for Integer
func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type for Integer
func (i *Integer) Type() Type {
	return IntType
}

// Float is a float64
type Float struct {
	Value float64
}

// String for Float
func (f *Float) String() string {
	return fmt.Sprintf("%f", f.Value)
}

// Type for Float
func (f *Float) Type() Type {
	return FloatType
}

// Boolean is a int64
type Boolean struct {
	Value bool
}

// String for Boolean
func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

// Type for Boolean
func (b *Boolean) Type() Type {
	return BooleanType
}

type Nil struct{}

// String for Boolean
func (n *Nil) String() string {
	return "nil"
}

// Type for Nil
func (n *Nil) Type() Type {
	return NilType
}
