package parser

import (
	"strings"

	"github.com/hazemKrimi/jack-compiler/internal/tokenizer"
)

func ParseTokens(tokens []tokenizer.Token) string {
	var output strings.Builder

	output.WriteString("<tokens>\n")

	for _, token := range tokens {
		switch token.Type {
		case tokenizer.SYMBOL:
			var value string

			switch token.Value {
			case "<":
				value = "&lt;"
			case ">":
				value = "&gt;"
			case "&":
				value = "&amp;"
			default:
				value = token.Value
			}

			output.WriteString("<symbol> " + value + " </symbol>\n")
		case tokenizer.KEYWORD:
			output.WriteString("<keyword> " + token.Value + " </keyword>\n")
		case tokenizer.IDENTIFIER:
			output.WriteString("<identifier> " + token.Value + " </identifier>\n")
		case tokenizer.INT_CONST:
			output.WriteString("<integerConstant> " + token.Value + " </integerConstant>\n")
		case tokenizer.STR_CONST:
			output.WriteString("<stringConstant> " + token.Value + " </stringConstant>\n")
		}
	}

	output.WriteString("</tokens>\n")

	return output.String()
}
