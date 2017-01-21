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
	Exp     // Exp ^
	Mod     // Mod %
	Bang    // Bang !
	Less    // Less <
	Greater // Greater >
	Inc     // Inc +=
	Dec     // Dec -=

	Equal    // Equal ==
	NotEqual // NotEqual !=

	Terminator // Terminator is the end of statement terminator

	LParen // LParen (
	RParen // RParen )

	LBrace // LBrace {
	RBrace // RBrace }

	LBracket // LBracket [
	RBracket // RBracket ]

	Dot      // Dot .
	Comma    // Comma ,
	Bar      // Bar |  - donotes function arg bar
	Continue // Continue - starts a single statment/line block

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
	case Equal:
		return "=="
	case NotEqual:
		return "!="
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Times:
		return "*"
	case Divide:
		return "/"
	case Exp:
		return "^"
	case Mod:
		return "%"
	case Inc:
		return "+="
	case Dec:
		return "-="
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
	case Continue:
		return ":"
	case Comma:
		return ","
	case Dot:
		return "."
	case Terminator:
		return "terminator"
	case EOF:
		return "EOF"
	case Identifier:
		return "identifier"
	case Int:
		return "int"
	case Float:
		return "float"
	case String:
		return "string"
	case Let:
		return "let"
	case If:
		return "if"
	case Else:
		return "else"
	case False:
		return "false"
	case True:
		return "true"
	case For:
		return "for"
	case Return:
		return "ret"
	default:
		return "unknown"
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
