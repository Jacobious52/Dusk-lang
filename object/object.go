package object

import (
	"fmt"
	"jacob/black/token"
)

// Type is a flag for the type of object
type Type int

const (
	// NilType nil
	NilType Type = iota
	// IntType int64
	IntType
	// FloatType float64
	FloatType
	// BooleanType bool
	BooleanType
	// StringType string
	StringType
	// ReturnType value
	ReturnType
	// ErrorType runtime error
	ErrorType
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
	case ReturnType:
		return "return_value"
	case ErrorType:
		return "error"
	default:
		return "unknown"
	}
}

// Object is the base interface of all objects
type Object interface {
	Type() Type
	String() string
	CanApply(token.Type, Type) bool
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

// CanApply for this type
func (i *Integer) CanApply(op token.Type, t Type) bool {
	switch t {
	case IntType, FloatType:
		return true
	default:
		if op == token.Equal || op == token.NotEqual {
			return true
		}
		return false
	}
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

// CanApply for this type
func (f *Float) CanApply(op token.Type, t Type) bool {
	switch t {
	case IntType, FloatType:
		return true
	default:
		return false
	}
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

// CanApply for this type
func (b *Boolean) CanApply(op token.Type, t Type) bool {
	switch op {
	case token.Equal, token.NotEqual:
		return true
	default:
		return false
	}
}

// Nil -  No value
type Nil struct{}

// String for Boolean
func (n *Nil) String() string {
	return "nil"
}

// Type for Nil
func (n *Nil) Type() Type {
	return NilType
}

// CanApply for this type
func (n *Nil) CanApply(op token.Type, t Type) bool {
	return false
}

// ReturnValue wrapper for a value returned
type ReturnValue struct {
	Value Object
}

// String for Return
func (r *ReturnValue) String() string {
	return r.Value.String()
}

// Type for Return
func (r *ReturnValue) Type() Type {
	return ReturnType
}

// CanApply for this type
func (r *ReturnValue) CanApply(op token.Type, t Type) bool {
	return r.Value.CanApply(op, t)
}

// Error for runrime error
type Error struct {
	Message string
	Pos     token.Position
}

// String for Error
func (e *Error) String() string {
	return fmt.Sprintf("%s: %s", e.Pos, e.Message)
}

// Type for Error
func (e *Error) Type() Type {
	return ErrorType
}

// CanApply for this type
func (e *Error) CanApply(op token.Type, t Type) bool {
	return false
}
