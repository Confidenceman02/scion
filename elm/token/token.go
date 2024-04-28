package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifies + literals
	IDENT   = "IDENT"
	INT     = "INT"
	LITERAL = "LITERAL"
)

var BINOP = [...]string{"+", "-", "/", "*", "=", ".", "<", ">", ":", "&", "|", "^", "?", "%", "!"}
