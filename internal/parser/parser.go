package parser

import (
	"github.com/hazemKrimi/jack-compiler/internal/tokenizer"
)

func ParseTokens(tokens []tokenizer.Token) string {
	output := "<tokens>\n"

	for _, token := range tokens {
		switch token.Type {
		case tokenizer.SYMBOL:
			{
				output += "<symbol>" + token.Value + "</symbol>\n"
			}
		case tokenizer.KEYWORD:
			{
				output += "<keyword>" + token.Value + "</keyword>\n"
			}

		}
	}

	output += "</tokens>\n"

	return output
}
