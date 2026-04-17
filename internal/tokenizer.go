package compiler

type TokenType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	INT_CONST
	STR_CONST
	IDENTIFIER
)

type Token struct {
	Value string
	Type TokenType
}

