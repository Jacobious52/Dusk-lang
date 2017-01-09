package token

// All valid token types
const (
	Illegal Type = (iota - 1) // Illegal type means it is a non recognised token
	EOF                       // EOF is the end of the file stream

	Identifier // Identifier any varible name

	Int    // Int literal type
	Float  // Double literal type
	String // Double literal type

	Assign  // Assign =
	Plus    // Plus +
	Minus   // Minus -
	Times   // Times *
	Divide  // Divide /
	Bang    // Bang !
	Less    // Less <
	Greater // Greater >

	Equal    // Equal ==
	NotEqual // NotEqual !=

	Comma      // Comma ,
	Terminator // Terminator is the end of statement terminator

	LParen // LParen (
	RParen // RParen )

	LBrace // LBrace {
	RBrace // RBrace }

	LBracket // LBracket [
	RBracket // RBracket ]

	Bar // Bar |  - donotes function arg bar

	Let    // Let keyword
	If     // If keyword
	Else   // Else keyword
	For    // For keyword
	Return // Ret keyword
	True   // True keyword
	False  // False keyword
)

// LookupLiteral returns string for type
// Only used in debugging and testing
func (t Type) String() string {
	switch t {
	case Assign:
		return "="
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Times:
		return "*"
	case Divide:
		return "/"
	case Bang:
		return "!"
	case Less:
		return "<"
	case Greater:
		return ">"
	case LBrace:
		return "{"
	case RBrace:
		return "}"
	case LParen:
		return "("
	case RParen:
		return ")"
	case LBracket:
		return "["
	case RBracket:
		return "]"
	case Bar:
		return "|"
	case Comma:
		return ","
	case Terminator:
		return "Terminator"
	case EOF:
		return "EOF"
	case Identifier:
		return "Identifier"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	case Let:
		return "Let"
	case If:
		return "If"
	case Else:
		return "Else"
	case False:
		return "False"
	case True:
		return "True"
	case For:
		return "For"
	case Return:
		return "Return"
	default:
		return "Unknown"
	}
}

// keywords maps the keyword to a Type
var keywords = map[string]Type{
	"let":   Let,
	"if":    If,
	"else":  Else,
	"false": False,
	"true":  True,
	"for":   For,
	"ret":   Return,
}
