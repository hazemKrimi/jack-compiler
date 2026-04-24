package parser

import (
	"errors"
	"slices"
	"strings"

	"github.com/hazemKrimi/jack-compiler/internal/tokenizer"
)

func WriteToken(output *strings.Builder, token tokenizer.Token, index *int) error {
	if _, err := output.WriteString("<" + token.XML + "> " + token.Value + " </" + token.XML + ">\n"); err != nil {
		return err
	}

	(*index)++

	return nil
}

func parseClassVarDec(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || !slices.Contains([]string{"static", "field"}, tokens[*index].Value) {
		return nil
	}

	output.WriteString("<classVarDec>\n")
	WriteToken(output, tokens[*index], index)

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid variable type name!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	WriteToken(output, tokens[*index], index)

	for tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		WriteToken(output, tokens[*index], index)

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid variable name!")
		}

		WriteToken(output, tokens[*index], index)
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	WriteToken(output, tokens[*index], index)
	output.WriteString("</classVarDec>\n")

	return parseClassVarDec(output, tokens, index)
}

func parseParameterList(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) || !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return nil
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		WriteToken(output, tokens[*index], index)

		return parseParameterList(output, tokens, index)
	}

	return nil
}

func parseVariableDeclaration(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "var" {
		return nil
	}

	output.WriteString("<varDec>\n")

	WriteToken(output, tokens[*index], index)

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid variable type name!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	WriteToken(output, tokens[*index], index)

	for tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		WriteToken(output, tokens[*index], index)

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid variable name!")
		}

		WriteToken(output, tokens[*index], index)
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	WriteToken(output, tokens[*index], index)
	output.WriteString("</varDec>\n")

	return parseVariableDeclaration(output, tokens, index)
}

func parseSubroutineCall(output *strings.Builder, tokens []tokenizer.Token, index *int) error {

	if tokens[*index].Value == "." {
		WriteToken(output, tokens[*index], index)

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid subroutine name!")
		}

		WriteToken(output, tokens[*index], index)
	}

	if tokens[*index].Value != "(" {
		return errors.New("Missing subroutine call opening parenthese!")
	}

	WriteToken(output, tokens[*index], index)

	output.WriteString("<expressionList>\n")

	if err := parseExpressionList(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</expressionList>\n")

	if tokens[*index].Value != ")" {
		return errors.New("Missing subroutine call closing parenthese!")
	}

	WriteToken(output, tokens[*index], index)

	return nil
}

func parseTerm(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	output.WriteString("<term>\n")

	if tokens[*index].Type == tokenizer.SYMBOL && slices.Contains([]string{"-", "~"}, tokens[*index].Value) {
		WriteToken(output, tokens[*index], index)

		if err := parseTerm(output, tokens, index); err != nil {
			return err
		}

		output.WriteString("</term>\n")
		return nil
	}

	if slices.Contains([]tokenizer.TokenType{tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) || slices.Contains([]string{"true", "false", "null", "this"}, tokens[*index].Value) {
		WriteToken(output, tokens[*index], index)
		output.WriteString("</term>\n")

		return nil
	}

	if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "(" {
		WriteToken(output, tokens[*index], index)

		if err := parseExpression(output, tokens, index); err != nil {
			return err
		}

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
			return errors.New("Invalid term!")
		}

		WriteToken(output, tokens[*index], index)
		output.WriteString("</term>\n")

		return nil
	}

	if tokens[*index].Type == tokenizer.IDENTIFIER {
		WriteToken(output, tokens[*index], index)

		if tokens[*index].Value == "[" {
			WriteToken(output, tokens[*index], index)

			if err := parseExpression(output, tokens, index); err != nil {
				return err
			}

			if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "]" {
				return errors.New("Invalid term!")
			}

			WriteToken(output, tokens[*index], index)
		} else if slices.Contains([]string{"(", "."}, tokens[*index].Value) {
			if err := parseSubroutineCall(output, tokens, index); err != nil {
				return err
			}
		}

		output.WriteString("</term>\n")
	}

	return nil
}

func parseExpression(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	output.WriteString("<expression>\n")

	if err := parseTerm(output, tokens, index); err != nil {
		return err
	}

	if slices.Contains([]string{"+", "-", "*", "/", "&amp;", "|", "&lt;", "&gt;", "="}, tokens[*index].Value) {
		WriteToken(output, tokens[*index], index)

		if err := parseTerm(output, tokens, index); err != nil {
			return err
		}
	}

	output.WriteString("</expression>\n")

	return nil
}

func parseExpressionList(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if slices.Contains([]tokenizer.TokenType{tokenizer.IDENTIFIER, tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) || slices.Contains([]string{"true", "false", "null", "this", "~", "-", "("}, tokens[*index].Value) {
		if err := parseExpression(output, tokens, index); err != nil {
			return err
		}

		if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
			WriteToken(output, tokens[*index], index)

			return parseExpressionList(output, tokens, index)
		}
	}

	return nil
}

func parseLetStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "let" {
		return errors.New("Invalid let statement!")
	}

	output.WriteString("<letStatement>\n")

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Value == "[" {
		WriteToken(output, tokens[*index], index)

		if err := parseExpression(output, tokens, index); err != nil {
			return err
		}

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "]" {
			return errors.New("Invalid expression!")
		}

		WriteToken(output, tokens[*index], index)
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "=" {
		return errors.New("Missing assignment!")
	}

	WriteToken(output, tokens[*index], index)

	if err := parseExpression(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	WriteToken(output, tokens[*index], index)
	output.WriteString("</letStatement>\n")

	return nil
}

func parseIfStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "if" {
		return errors.New("Invalid if statement!")
	}

	output.WriteString("<ifStatement>\n")

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing if statement opening parenthese!")
	}

	WriteToken(output, tokens[*index], index)

	if err := parseExpression(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing if statement closing parenthese!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing if statement opening curly brace!")
	}

	WriteToken(output, tokens[*index], index)

	output.WriteString("<statements>\n")

	if err := parseStatements(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</statements>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing if statement closing curly brace!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type == tokenizer.KEYWORD && tokens[*index].Value == "else" {
		WriteToken(output, tokens[*index], index)

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
			return errors.New("Missing if statement opening curly brace!")
		}

		WriteToken(output, tokens[*index], index)
		output.WriteString("<statements>\n")

		if err := parseStatements(output, tokens, index); err != nil {
			return err
		}

		output.WriteString("</statements>\n")

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
			return errors.New("Missing if statement closing curly brace!")
		}

		WriteToken(output, tokens[*index], index)
	}

	output.WriteString("</ifStatement>\n")

	return nil
}

func parseWhileStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "while" {
		return errors.New("Invalid while statement!")
	}

	output.WriteString("<whileStatement>\n")

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing while statement opening parenthese!")
	}

	WriteToken(output, tokens[*index], index)

	if err := parseExpression(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing while statement closing parenthese!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing while statement opening curly brace!")
	}

	WriteToken(output, tokens[*index], index)

	output.WriteString("<statements>\n")

	if err := parseStatements(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</statements>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing while statement closing curly brace!")
	}

	WriteToken(output, tokens[*index], index)

	output.WriteString("</whileStatement>\n")

	return nil
}

func parseDoStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "do" {
		return errors.New("Invalid do statement!")
	}

	output.WriteString("<doStatement>\n")

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	WriteToken(output, tokens[*index], index)

	if err := parseSubroutineCall(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	WriteToken(output, tokens[*index], index)
	output.WriteString("</doStatement>\n")

	return nil
}

func parseReturnStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "return" {
		return errors.New("Invalid return statement!")
	}

	output.WriteString("<returnStatement>\n")

	WriteToken(output, tokens[*index], index)

	if slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER, tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) {
		if err := parseExpression(output, tokens, index); err != nil {
			return err
		}
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	WriteToken(output, tokens[*index], index)
	output.WriteString("</returnStatement>\n")

	return nil
}

func parseStatements(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD {
		return nil
	}

	switch tokens[*index].Value {
	case "let":
		if err := parseLetStatement(output, tokens, index); err != nil {
			return err
		}
	case "if":
		if err := parseIfStatement(output, tokens, index); err != nil {
			return err
		}
	case "while":
		if err := parseWhileStatement(output, tokens, index); err != nil {
			return err
		}
	case "do":
		if err := parseDoStatement(output, tokens, index); err != nil {
			return err
		}
	case "return":
		if err := parseReturnStatement(output, tokens, index); err != nil {
			return err
		}
	default:
		return errors.New("Invalid statement!")
	}

	return parseStatements(output, tokens, index)
}

func parseSubroutineBody(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type == tokenizer.KEYWORD && tokens[*index].Value == "var" {
		if err := parseVariableDeclaration(output, tokens, index); err != nil {
			return err
		}
	}

	output.WriteString("<statements>\n")

	if err := parseStatements(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</statements>\n")

	return nil
}

func parseSubroutineDeclaration(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || !slices.Contains([]string{"constructor", "method", "function"}, tokens[*index].Value) {
		return nil
	}

	output.WriteString("<subroutineDec>\n")

	WriteToken(output, tokens[*index], index)

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"void", "int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid subroutine return type!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid subroutine name!")
	}

	WriteToken(output, tokens[*index], index)

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing subroutine opening parenthese!")
	}

	WriteToken(output, tokens[*index], index)
	output.WriteString("<parameterList>\n")

	if err := parseParameterList(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</parameterList>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing subroutine closing parenthese!")
	}

	WriteToken(output, tokens[*index], index)
	output.WriteString("<subroutineBody>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing subroutine opening curly brace!")
	}

	WriteToken(output, tokens[*index], index)

	if err := parseSubroutineBody(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing subroutine closing curly brace!")
	}

	WriteToken(output, tokens[*index], index)
	output.WriteString("</subroutineBody>\n")
	output.WriteString("</subroutineDec>\n")

	return parseSubroutineDeclaration(output, tokens, index)
}

func parseClass(output *strings.Builder, tokens []tokenizer.Token) error {
	index := 0

	output.WriteString("<class>\n")

	if tokens[index].Type != tokenizer.KEYWORD || tokens[index].Value != "class" {
		return errors.New("Jack file must contain one class!")
	}

	WriteToken(output, tokens[index], &index)

	if tokens[index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid class name!")
	}

	WriteToken(output, tokens[index], &index)

	if tokens[index].Type != tokenizer.SYMBOL || tokens[index].Value != "{" {
		return errors.New("Missing class opening curly brace!")
	}

	WriteToken(output, tokens[index], &index)

	if err := parseClassVarDec(output, tokens, &index); err != nil {
		return err
	}

	if err := parseSubroutineDeclaration(output, tokens, &index); err != nil {
		return err
	}

	if tokens[index].Type != tokenizer.SYMBOL || tokens[index].Value != "}" {
		return errors.New("Missing class closing curly brace!")
	}

	WriteToken(output, tokens[index], &index)
	output.WriteString("</class>\n")

	return nil
}

func ParseTokens(tokens []tokenizer.Token) (string, error) {
	var output strings.Builder

	if err := parseClass(&output, tokens); err != nil {
		return "", err
	}

	return output.String(), nil
}
