package tokenizer

import (
	"bufio"
	"regexp"
	"slices"
)

type TokenType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	INT_CONST
	STR_CONST
	IDENTIFIER
)

var KEYWORDS = []string{
	"class",
	"constructor",
	"function",
	"method",
	"field",
	"static",
	"var",
	"int",
	"read",
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

var SYMBOLS = []string{
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

func ExtractTokens(tokens *[]Token, reader *bufio.Reader) error {
	read := ""
	buf := []byte{}

	for {
		_, err := reader.Read(buf)
		read += string(buf)

		if err != nil {
			return err
		}

		if match, _ := regexp.MatchString("^[[:space:]]$", read); match {
			continue
		}

		if read == "/" {
			next, err := reader.ReadByte()

			if err != nil {
				return err
			}

			if string(next) == "/" || string(next) == "*" {
				_, err := reader.ReadBytes('/')

				if err != nil {
					return err
				}
			} else {
				*tokens = append(*tokens, Token{Value: read, Type: SYMBOL})
				read = ""

				err := reader.UnreadByte()

				if err != nil {
					return err
				}

				continue
			}
		}

		if slices.Contains(SYMBOLS, read) {
			*tokens = append(*tokens, Token{Value: read, Type: SYMBOL})
			read = ""

			continue
		}

		if slices.Contains(KEYWORDS, read) {
			*tokens = append(*tokens, Token{Value: read, Type: KEYWORD})
			read = ""

			continue
		}
	}

	return nil
}
