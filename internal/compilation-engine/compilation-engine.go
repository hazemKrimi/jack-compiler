package parser

import (
	"errors"
	"slices"
	"strings"

	"github.com/hazemKrimi/jack-compiler/internal/tokenizer"
)

func compileClassVarDec(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || !slices.Contains([]string{"static", "field"}, tokens[*index].Value) {
		return nil
	}

	output.WriteString("<classVarDec>\n")
	output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	*(index)++

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid variable type name!")
	}

	if tokens[*index].Type == tokenizer.KEYWORD {
		output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	} else {
		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	}

	*(index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	*(index)++

	for tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
		(*index)++

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid variable name!")
		}

		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
		(*index)++
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	(*index)++

	output.WriteString("</classVarDec>\n")

	return compileClassVarDec(output, tokens, index)
}

func compileParameterList(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) || !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return nil
	}

	if tokens[*index].Type == tokenizer.KEYWORD {
		output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	} else {
		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	}

	*(index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	*(index)++

	if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
		*(index)++

		return compileParameterList(output, tokens, index)
	}

	return nil
}

func compileVariableDeclaration(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "var" {
		return nil
	}

	output.WriteString("<varDec>\n")

	output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	*(index)++

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid variable type name!")
	}

	if tokens[*index].Type == tokenizer.KEYWORD {
		output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	} else {
		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	}

	*(index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	*(index)++

	for tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
		(*index)++

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid variable name!")
		}

		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
		(*index)++
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	(*index)++

	output.WriteString("</varDec>\n")

	return compileVariableDeclaration(output, tokens, index)
}

func compileExpression(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER, tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) && !slices.Contains([]string{"true", "false", "null", "this"}, tokens[*index].Value) {
		return errors.New("Invalid expression!")
	}

	output.WriteString("<expression>\n")
	output.WriteString("<term>\n")

	switch tokens[*index].Type {
	// case tokenizer.SYMBOL:
	// 	var value string
	//
	// 	switch tokens[*index].Value {
	// 	case "<":
	// 		value = "&lt;"
	// 	case ">":
	// 		value = "&gt;"
	// 	case "&":
	// 		value = "&amp;"
	// 	default:
	// 		value = tokens[*index].Value
	// 	}
	//
	// 	output.WriteString("<symbol> " + value + " </symbol>\n")
	case tokenizer.KEYWORD:
		output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	case tokenizer.IDENTIFIER:
		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	case tokenizer.INT_CONST:
		output.WriteString("<integerConstant> " + tokens[*index].Value + " </integerConstant>\n")
	case tokenizer.STR_CONST:
		output.WriteString("<stringConstant> " + tokens[*index].Value + " </stringConstant>\n")
	}

	output.WriteString("</term>\n")
	output.WriteString("</expression>\n")
	*(index)++

	return nil
}

func compileExpressionList(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER, tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) || slices.Contains([]string{"true", "false", "null", "this"}, tokens[*index].Value) {
		if err := compileExpression(output, tokens, index); err != nil {
			return err
		}

		if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
			output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
			(*index)++

			return compileExpressionList(output, tokens, index)
		}
	}

	return nil
}

func compileLetStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "let" {
		return errors.New("Invalid let statement!")
	}

	output.WriteString("<letStatement>\n")

	output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "=" {
		return errors.New("Missing assignment!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	(*index)++

	if err := compileExpression(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	(*index)++

	output.WriteString("</letStatement>\n")

	return nil
}

func compileIfStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "if" {
		return errors.New("Invalid if statement!")
	}

	output.WriteString("<ifStatement>\n")

	output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing if statement opening parenthese!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	if err := compileExpression(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing if statement closing parenthese!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing if statement opening curly brace!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	output.WriteString("<statements>\n")

	if err := compileStatements(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</statements>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing if statement closing curly brace!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	if tokens[*index].Type == tokenizer.KEYWORD && tokens[*index].Value == "else" {
		output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
		*(index)++

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
			return errors.New("Missing if statement opening curly brace!")
		}

		output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
		*(index)++

		output.WriteString("<statements>\n")

		if err := compileStatements(output, tokens, index); err != nil {
			return err
		}

		output.WriteString("</statements>\n")

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
			return errors.New("Missing if statement closing curly brace!")
		}

		output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
		*(index)++
	}

	output.WriteString("</ifStatement>\n")

	return nil
}

func compileWhileStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "while" {
		return errors.New("Invalid while statement!")
	}

	output.WriteString("<whileStatement>\n")

	output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing while statement opening parenthese!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	if err := compileExpression(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing while statement closing parenthese!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing while statement opening curly brace!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	output.WriteString("<statements>\n")

	if err := compileStatements(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</statements>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing while statement closing curly brace!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	output.WriteString("</whileStatement>\n")

	return nil
}

func compileDoStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "do" {
		return errors.New("Invalid do statement!")
	}

	output.WriteString("<doStatement>\n")

	output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	*(index)++

	if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "." {
		output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
		*(index)++

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid variable name!")
		}

		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
		*(index)++
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing subroutine call opening parenthese!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	output.WriteString("<expressionList>\n")

	if err := compileExpressionList(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</expressionList>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing subroutine call closing parenthese!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	(*index)++

	output.WriteString("</doStatement>\n")

	return nil
}

func compileReturnStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "return" {
		return errors.New("Invalid return statement!")
	}

	output.WriteString("<returnStatement>\n")

	output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	*(index)++

	if slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER, tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) {
		if err := compileExpression(output, tokens, index); err != nil {
			return err
		}
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	(*index)++

	output.WriteString("</returnStatement>\n")

	return nil
}

func compileStatements(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD {
		return nil
	}

	switch tokens[*index].Value {
	case "let":
		if err := compileLetStatement(output, tokens, index); err != nil {
			return err
		}
	case "if":
		if err := compileIfStatement(output, tokens, index); err != nil {
			return err
		}
	case "while":
		if err := compileWhileStatement(output, tokens, index); err != nil {
			return err
		}
	case "do":
		if err := compileDoStatement(output, tokens, index); err != nil {
			return err
		}
	case "return":
		if err := compileReturnStatement(output, tokens, index); err != nil {
			return err
		}
	default:
		return errors.New("Invalid statement!")
	}

	return compileStatements(output, tokens, index)
}

func compileSubroutineBody(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type == tokenizer.KEYWORD && tokens[*index].Value == "var" {
		if err := compileVariableDeclaration(output, tokens, index); err != nil {
			return err
		}
	}

	output.WriteString("<statements>\n")

	if err := compileStatements(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</statements>\n")

	return nil
}

func compileSubroutineDeclaration(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || !slices.Contains([]string{"constructor", "method", "function"}, tokens[*index].Value) {
		return nil
	}

	output.WriteString("<subroutineDec>\n")

	output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	*(index)++

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"void", "int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid subroutine return type!")
	}

	if tokens[*index].Type == tokenizer.KEYWORD {
		output.WriteString("<keyword> " + tokens[*index].Value + " </keyword>\n")
	} else {
		output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	}

	*(index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid subroutine name!")
	}

	output.WriteString("<identifier> " + tokens[*index].Value + " </identifier>\n")
	*(index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing subroutine opening parenthese!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	output.WriteString("<parameterList>\n")

	if err := compileParameterList(output, tokens, index); err != nil {
		return err
	}

	output.WriteString("</parameterList>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing subroutine closing parenthese!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	output.WriteString("<subroutineBody>\n")

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing subroutine opening curly brace!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	if err := compileSubroutineBody(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing subroutine closing curly brace!")
	}

	output.WriteString("<symbol> " + tokens[*index].Value + " </symbol>\n")
	*(index)++

	output.WriteString("</subroutineBody>\n")
	output.WriteString("</subroutineDec>\n")

	return compileSubroutineDeclaration(output, tokens, index)
}

func compileClass(output *strings.Builder, tokens []tokenizer.Token) error {
	index := 0

	output.WriteString("<class>\n")

	if tokens[index].Type != tokenizer.KEYWORD || tokens[index].Value != "class" {
		return errors.New("Jack file must contain one class!")
	}

	output.WriteString("<keyword> " + tokens[index].Value + " </keyword>\n")
	index++

	if tokens[index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid class name!")
	}

	output.WriteString("<identifier> " + tokens[index].Value + " </identifier>\n")
	index++

	if tokens[index].Type != tokenizer.SYMBOL || tokens[index].Value != "{" {
		return errors.New("Missing class opening curly brace!")
	}

	output.WriteString("<symbol> " + tokens[index].Value + " </symbol>\n")
	index++

	if err := compileClassVarDec(output, tokens, &index); err != nil {
		return err
	}

	if err := compileSubroutineDeclaration(output, tokens, &index); err != nil {
		return err
	}

	if tokens[index].Type != tokenizer.SYMBOL || tokens[index].Value != "}" {
		return errors.New("Missing class closing curly brace!")
	}

	output.WriteString("<symbol> " + tokens[index].Value + " </symbol>\n")
	output.WriteString("</class>\n")

	return nil
}

func ParseTokens(tokens []tokenizer.Token) (string, error) {
	var output strings.Builder

	if err := compileClass(&output, tokens); err != nil {
		return "", err
	}

	return output.String(), nil
}
