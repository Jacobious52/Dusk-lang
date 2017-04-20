package eval

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"jacob/dusk/object"
	"jacob/dusk/token"
	"os"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len":     &object.Builtin{Fn: length},
	"first":   &object.Builtin{Fn: first},
	"last":    &object.Builtin{Fn: last},
	"rest":    &object.Builtin{Fn: rest},
	"lead":    &object.Builtin{Fn: lead},
	"push":    &object.Builtin{Fn: push},
	"pop":     &object.Builtin{Fn: pop},
	"alloc":   &object.Builtin{Fn: alloc},
	"set":     &object.Builtin{Fn: set},
	"join":    &object.Builtin{Fn: join},
	"split":   &object.Builtin{Fn: split},
	"println": &object.Builtin{Fn: println},
	"print":   &object.Builtin{Fn: print},
	"readln":  &object.Builtin{Fn: readln},
	"read":    &object.Builtin{Fn: read},
	"readc":   &object.Builtin{Fn: readc},
	"readall": &object.Builtin{Fn: readall},
	"atoi":    &object.Builtin{Fn: atoi},
	"itoa":    &object.Builtin{Fn: itoa},
	"in":      &object.Builtin{Fn: in},
	"out":     &object.Builtin{Fn: out},
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

func lead(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '1'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		l := len(arg.Value)
		if l > 0 {
			newStr := make([]byte, l-1, l-1)
			copy(newStr, arg.Value[:l-1])
			return &object.String{Value: string(newStr)}
		}
		return ConstNil
	case *object.Array:
		l := len(arg.Elements)
		if l > 0 {
			newElems := make([]object.Object, l-1, l-1)
			copy(newElems, arg.Elements[:l-1])
			return &object.Array{Elements: newElems}
		}
		return ConstNil
	default:
		return newError(token.Position{}, "argument to 'lead' not supported, got '%s'", args[0].Type())
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

func pop(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '1'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		l := len(arg.Value)
		if l > 0 {
			p := arg.Value[l-1]
			arg.Value = arg.Value[:l-1]
			return &object.String{Value: string(p)}
		}
		return ConstNil
	case *object.Array:
		l := len(arg.Elements)
		if l > 0 {
			p := arg.Elements[l-1]
			arg.Elements = arg.Elements[:l-1]
			return p
		}
		return ConstNil
	default:
		return newError(token.Position{}, "argument to 'pop' not supported, got '%s'", args[0].Type())
	}
}

func alloc(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '2'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		if arg.Value >= 0 {
			newArr := make([]object.Object, arg.Value, arg.Value)
			for i := range newArr {
				newArr[i] = args[1]
			}
			return &object.Array{Elements: newArr}
		}
		return ConstNil
	default:
		return newError(token.Position{}, "argument to 'alloc' not supported, got '%s'", arg.Type())
	}
}

func set(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '3'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if i, ok := args[1].(*object.Integer); ok {
			arg.Elements[i.Value] = args[2]
			return ConstNil
		}
		return newError(token.Position{}, "second argument to 'set' not supported, got '%s'", args[1].Type())
	default:
		return newError(token.Position{}, "argument to 'set' not supported, got '%s'", arg.Type())
	}
}

func join(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '2'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if s, ok := args[1].(*object.String); ok {
			parts := make([]string, len(arg.Elements), len(arg.Elements))
			for i := range arg.Elements {
				parts[i] = arg.Elements[i].String()
			}
			return &object.String{Value: strings.Join(parts, s.Value)}
		}
		return newError(token.Position{}, "second argument to 'join' not supported, got '%s'", args[1].Type())
	default:
		return newError(token.Position{}, "argument to 'join' not supported, got '%s'", arg.Type())
	}
}

func split(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '2'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		if s, ok := args[1].(*object.String); ok {

			parts := strings.Split(arg.Value, s.Value)

			elems := make([]object.Object, len(parts), len(parts))
			for i := range parts {
				elems[i] = &object.String{Value: parts[i]}
			}
			return &object.Array{Elements: elems}
		}
		return newError(token.Position{}, "second argument to 'split' not supported, got '%s'", args[1].Type())
	default:
		return newError(token.Position{}, "argument to 'split' not supported, got '%s'", arg.Type())
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

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return &object.String{Value: scanner.Text()}
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

	reader := bufio.NewReader(os.Stdin)
	c, e := reader.ReadByte()
	if e != nil {
		return ConstNil
	}

	return &object.String{Value: string(c)}
}

func readall(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError(token.Position{}, "readln does not take any arguments. given '%d'", len(args))
	}

	s, _ := ioutil.ReadAll(os.Stdin)

	return &object.String{Value: string(s)}
}

func in(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "in takes one arguments. given '%d'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		f, err := os.Open(arg.Value)
		if err != nil {
			return newError(token.Position{}, err.Error())
		}
		defer f.Close()
		s, _ := ioutil.ReadAll(f)
		return &object.String{Value: string(s)}
	default:
		return newError(token.Position{}, "argument to 'in' not supported, got '%s'", args[0].Type())
	}
}

func out(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError(token.Position{}, "out takes one arguments. given '%d'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		f, err := os.Create(arg.Value)
		if err != nil {
			return newError(token.Position{}, err.Error())
		}
		defer f.Close()
		switch str := args[1].(type) {
		case *object.String:
			f.WriteString(str.Value)
			return nil
		default:
			return newError(token.Position{}, "argument to 'out' not supported, got '%s'", args[0].Type())
		}
	default:
		return newError(token.Position{}, "argument to 'out' not supported, got '%s'", args[0].Type())
	}
}

func atoi(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '1'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		l := len(arg.Value)
		if l == 1 {
			return &object.Integer{Value: int64(arg.Value[0])}
		}
		return newError(token.Position{}, "argument to 'atoi must be string with length of 1. Got '%d'", l)
	default:
		return newError(token.Position{}, "argument to 'atoi' not supported, got '%s'", args[0].Type())
	}
}

func itoa(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError(token.Position{}, "wrong number of arguments. got '%d', expected '1'", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		if arg.Value >= 0 && arg.Value < 256 {
			return &object.String{Value: string(byte(arg.Value))}
		}
		return newError(token.Position{}, "argument to 'atoi must be between 0 and 256 Got '%d'", arg.Value)
	default:
		return newError(token.Position{}, "argument to 'atoi' not supported, got '%s'", args[0].Type())
	}
}
