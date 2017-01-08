package token

// All valid token types
const (
	// Illegal type means it is a non recognised token
	Illegal Type = (iota - 1)
	// EOF is the end of the file stream
	EOF

	// Identifier any varible name
	Identifier

	// Int literal type
	Int
	// Double literal type
	Double
	// String literal type
	String

	// Assign =
	Assign
	// Plus +
	Plus
	// Minus -
	Minus
	// Times *
	Times
	// Divide /
	Divide
	// Bang !
	Bang
	// Less <
	Less
	// Greater >
	Greater

	// Equal ==
	Equal
	// NotEqual !=
	NotEqual

	// Comma ,
	Comma
	// Terminator is the end of statement terminator
	Terminator

	// LParen (
	LParen
	// RParen )
	RParen

	// LBrace {
	LBrace
	// RBrace }
	RBrace

	// LBracket [
	LBracket
	// RBracket ]
	RBracket

	// Bar |  - donotes function arg bar
	Bar

	// Let keyword
	Let
	// If keyword
	If
	// Else keyword
	Else
	// For keyword
	For
	// Return keyword
	Return
	// True keyword
	True
	// False keyword
	False
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
	case Double:
		return "Double"
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
