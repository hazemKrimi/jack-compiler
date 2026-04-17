package tokenizer

import (
	"fmt"
	"regexp"
)

type TokenType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	INT_CONST
	STR_CONST
	IDENTIFIER
)

var KEYWORDS = [...]string{
	"class",
	"constructor",
	"function",
	"method",
	"field",
	"static",
	"var",
	"int",
	"char",
	"boolean",
	"void",
	"true",
	"false",
	"null",
	"this",
	"let",
	"do",
	"if",
	"else",
	"while",
	"return",
}

var SYMBOLS = [...]string{
	"{",
	"}",
	"(",
	")",
	"[",
	"]",
	".",
	",",
	";",
	"+",
	"-",
	"*",
	"/",
	"&",
	"|",
	"<",
	">",
	"=",
	"~",
}

type Token struct {
	Value string
	Type  TokenType
}

func ExtractTokens(tokens *[]Token, source []byte) error {
	i := 0

	for i < len(source) {
		t := string(source[i])

		if match, _ := regexp.MatchString("^[[:space:]]$", t); match {
			i++
		} else {
			fmt.Println(t)
			i++
		}
	}

	return nil
}
