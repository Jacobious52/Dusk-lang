package eval

import (
	"fmt"
	"jacob/dusk/object"
	"jacob/dusk/token"
)

var builtins = map[string]*object.Builtin{
	"len":     &object.Builtin{Fn: length},
	"first":   &object.Builtin{Fn: first},
	"last":    &object.Builtin{Fn: last},
	"rest":    &object.Builtin{Fn: rest},
	"push":    &object.Builtin{Fn: push},
	"println": &object.Builtin{Fn: println},
	"print":   &object.Builtin{Fn: print},
	"readln":  &object.Builtin{Fn: readln},
	"read":    &object.Builtin{Fn: read},
	"readc":   &object.Builtin{Fn: readc},
}

func length(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '1'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError(token.Position{}, "argument to 'len' not supported, got '%s'", args[0].Type())
	}
}

func first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '1'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		if len(arg.Value) > 0 {
			return &object.String{Value: string(arg.Value[0])}
		}
		return ConstNil
	case *object.Array:
		if len(arg.Elements) > 0 {
			return arg.Elements[0]
		}
		return ConstNil
	default:
		return newError(token.Position{}, "argument to 'first' not supported, got '%s'", args[0].Type())
	}
}

func last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '1'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		if len(arg.Value) > 0 {
			return &object.String{Value: string(arg.Value[len(arg.Value)-1])}
		}
		return ConstNil
	case *object.Array:
		if len(arg.Elements) > 0 {
			return arg.Elements[len(arg.Elements)-1]
		}
		return ConstNil
	default:
		return newError(token.Position{}, "argument to 'last' not supported, got '%s'", args[0].Type())
	}
}

func rest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '1'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		l := len(arg.Value)
		if l > 0 {
			newStr := make([]byte, l-1, l-1)
			copy(newStr, arg.Value[1:l])
			return &object.String{Value: string(newStr)}
		}
		return ConstNil
	case *object.Array:
		l := len(arg.Elements)
		if l > 0 {
			newElems := make([]object.Object, l-1, l-1)
			copy(newElems, arg.Elements[1:l])
			return &object.Array{Elements: newElems}
		}
		return ConstNil
	default:
		return newError(token.Position{}, "argument to 'rest' not supported, got '%s'", args[0].Type())
	}
}

func push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '2'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		if p, ok := args[1].(*object.String); ok {
			str := arg.Value + p.Value
			return &object.String{Value: str}
		}
		return newError(token.Position{}, "cannot push '%s' to string", args[1].Type())
	case *object.Array:
		l := len(arg.Elements)
		newElems := make([]object.Object, l+1, l+1)
		copy(newElems, arg.Elements)
		newElems[l] = args[1]
		return &object.Array{Elements: newElems}
	default:
		return newError(token.Position{}, "argument to 'push' not supported, got '%s'", args[0].Type())
	}
}

func println(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg)
	}
	return ConstNil
}

func print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg)
	}
	return ConstNil
}

func readln(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError(token.Position{}, "readln does not take any arguments. given '%d'", len(args))
	}

	s := ""
	fmt.Scanln(&s)

	return &object.String{Value: s}
}

func read(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError(token.Position{}, "readln does not take any arguments. given '%d'", len(args))
	}

	s := ""
	fmt.Scan(&s)

	return &object.String{Value: s}
}

func readc(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError(token.Position{}, "readln does not take any arguments. given '%d'", len(args))
	}

	var c byte
	fmt.Scan(c)

	return &object.String{Value: string(c)}
}
