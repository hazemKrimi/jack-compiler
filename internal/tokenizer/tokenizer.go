package tokenizer

import (
	"bufio"
	"regexp"
	"slices"
	"strings"

	"github.com/hazemKrimi/jack-compiler/internal/utils"
)

type TokenType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STR_CONST
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
	"char",
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
	XML   string
}

func NewToken(value string, tokenType TokenType) Token {
	switch tokenType {
	case SYMBOL:
		return Token{Value: value, Type: tokenType, XML: "symbol"}
	case KEYWORD:
		return Token{Value: value, Type: tokenType, XML: "keyword"}
	case IDENTIFIER:
		return Token{Value: value, Type: tokenType, XML: "identifier"}
	case INT_CONST:
		return Token{Value: value, Type: tokenType, XML: "integerConstant"}
	case STR_CONST:
		return Token{Value: value, Type: tokenType, XML: "stringConstant"}
	default:
		return Token{Value: value, Type: tokenType, XML: "token"}
	}
}

func isDigit(text string) (bool, error) {
	if match, err := regexp.MatchString("^[[:digit:]]$", text); err != nil {
		return false, err
	} else if match {
		return true, nil
	}

	return false, nil
}

func isWhiteSpace(text string) (bool, error) {
	if match, err := regexp.MatchString("^[[:space:]]$", text); err != nil {
		return false, err
	} else if match {
		return true, nil
	}

	return false, nil
}

func isComment(text string, reader *bufio.Reader) (bool, error) {
	if text == "/" {
		n, err := reader.ReadByte()
		next := string(n)

		if err != nil {
			return false, err
		}

		switch next {
		case "/":
			_, err := reader.ReadBytes('\n')

			if err != nil {
				return false, err
			}

			return true, nil
		case "*":
			_, err := reader.ReadBytes('*')

			if err != nil {
				return false, err
			}

			p, err := reader.Peek(1)

			if err != nil {
				return false, err
			}

			for string(p) != "/" {
				_, err := reader.ReadBytes('*')

				if err != nil {
					return false, err
				}

				p, err = reader.Peek(1)
			}

			_, err = reader.ReadBytes('/')

			if err != nil {
				return false, err
			}

			return true, nil
		default:
			err = reader.UnreadByte()

			if err != nil {
				return false, err
			}

			return false, nil
		}
	}

	return false, nil
}

func ExtractTokens(tokens *[]Token, reader *bufio.Reader) error {
	buf := make([]byte, 0)

	for {
		b, err := reader.ReadByte()

		if isErr, isEOF := utils.CheckReaderError(err); isEOF {
			break
		} else if isErr {
			return err
		}

		text := string(b)
		comment, err := isComment(text, reader)

		if isErr, isEOF := utils.CheckReaderError(err); isEOF {
			break
		} else if isErr {
			return err
		} else if comment {
			continue
		}

		whitespace, err := isWhiteSpace(text)

		if isErr, isEOF := utils.CheckReaderError(err); isEOF {
			break
		} else if isErr {
			return err
		} else if whitespace {
			if len(buf) > 0 {

				read := string(buf)

				if slices.Contains(KEYWORDS, read) {
					*tokens = append(*tokens, NewToken(read, KEYWORD))
				} else {
					*tokens = append(*tokens, NewToken(read, IDENTIFIER))
				}

				buf = nil
			}

			continue
		}

		digit, err := isDigit(text)

		if isErr, isEOF := utils.CheckReaderError(err); isEOF {
			break
		} else if isErr {
			return err
		} else if digit {
			if len(buf) == 0 {
				var integerConstant strings.Builder

				integerConstant.WriteString(text)

				for {
					b, err := reader.ReadByte()

					if isErr, isEOF := utils.CheckReaderError(err); isEOF {
						break
					} else if isErr {
						return err
					}

					anotherDigit, err := isDigit(string(b))

					if isErr, isEOF := utils.CheckReaderError(err); isEOF {
						break
					} else if isErr {
						return err
					}

					if !anotherDigit {
						err := reader.UnreadByte()

						if isErr, isEOF := utils.CheckReaderError(err); isEOF {
							break
						} else if isErr {
							return err
						}

						break
					} else {
						integerConstant.WriteString(string(b))
					}
				}

				*tokens = append(*tokens, NewToken(integerConstant.String(), INT_CONST))
				continue
			}
		}

		if text == "\"" {
			b, err := reader.ReadBytes('"')

			if isErr, isEOF := utils.CheckReaderError(err); isEOF {
				break
			} else if isErr {
				return err
			}

			*tokens = append(*tokens, NewToken(string(b[:len(b)-1]), STR_CONST))
			continue
		}

		if slices.Contains(SYMBOLS, text) {
			if len(buf) > 0 {
				read := string(buf)

				if slices.Contains(KEYWORDS, read) {
					*tokens = append(*tokens, NewToken(read, KEYWORD))
				} else {
					*tokens = append(*tokens, NewToken(read, IDENTIFIER))
				}

				buf = nil
			}

			*tokens = append(*tokens, NewToken(text, SYMBOL))

			continue
		}

		buf = append(buf, b)
	}

	return nil
}
