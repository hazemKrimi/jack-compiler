package engine

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/hazemKrimi/jack-compiler/internal/code"
	"github.com/hazemKrimi/jack-compiler/internal/symbol-table"
	"github.com/hazemKrimi/jack-compiler/internal/tokenizer"
)

var className string
var classSymbolTable, subroutineSymbolTable map[string]table.Variable
var ifLabelIndex = 0
var whileLabelIndex = 0

func segmentFromVariableKind(kind table.VariableKind) code.Segment {
	switch kind {
	case table.FIELD:
		return code.THIS
	case table.VAR:
		return code.LOCAL
	case table.STATIC:
		return code.STATIC
	case table.ARG:
		return code.ARGUMENT
	}

	return -1
}

func isMethodName(tokens []tokenizer.Token, identifier string) bool {
	for index, token := range tokens {
		if token.Value == identifier && index >= 2 && tokens[index-2].Value == "method" {
			return true
		}
	}

	return false
}

func compileTerm(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type == tokenizer.SYMBOL && slices.Contains([]string{"-", "~"}, tokens[*index].Value) {
		op := tokens[*index].Value

		(*index)++

		if err := compileTerm(output, tokens, index); err != nil {
			return err
		}

		switch op {
		case "-":
			code.WriteArithmeticLogical(output, code.NEG)
		case "~":
			code.WriteArithmeticLogical(output, code.NOT)
		}

		return nil
	}

	if slices.Contains([]tokenizer.TokenType{tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) || slices.Contains([]string{"true", "false", "null", "this"}, tokens[*index].Value) {
		if tokens[*index].Type == tokenizer.INT_CONST {
			number, err := strconv.ParseInt(tokens[*index].Value, 10, 32)

			if err != nil {
				return err
			}

			code.WritePush(output, code.CONSTANT, int(number))
		}

		if tokens[*index].Value == "this" {
			variable, found := table.GetVariable([]*map[string]table.Variable{&subroutineSymbolTable}, tokens[*index].Value)

			if found {
				code.WritePush(output, segmentFromVariableKind(variable.Kind), variable.Count)
			} else {
				code.WritePush(output, code.POINTER, 0)
			}
		}

		if tokens[*index].Value == "true" {
			code.WritePush(output, code.CONSTANT, 1)
			code.WriteArithmeticLogical(output, code.NEG)
		}

		if slices.Contains([]string{"false", "null"}, tokens[*index].Value) {
			code.WritePush(output, code.CONSTANT, 0)
		}

		(*index)++

		return nil
	}

	if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "(" {
		(*index)++

		if err := compileExpression(output, tokens, index); err != nil {
			return err
		}

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
			return errors.New("Invalid term!")
		}

		(*index)++

		return nil
	}

	if tokens[*index].Type == tokenizer.IDENTIFIER {
		variable, found := table.GetVariable([]*map[string]table.Variable{&subroutineSymbolTable, &classSymbolTable}, tokens[*index].Value)

		if !found && !slices.Contains([]string{"(", "."}, tokens[(*index)+1].Value) {
			return errors.New("Invalid identifier!")
		}

		var identifier string

		if found && !slices.Contains([]string{"(", "."}, tokens[(*index)+1].Value) {
			code.WritePush(output, segmentFromVariableKind(variable.Kind), variable.Count)
		} else {
			identifier = tokens[*index].Value
		}

		(*index)++

		if tokens[*index].Value == "[" {
			(*index)++

			if err := compileExpression(output, tokens, index); err != nil {
				return err
			}

			if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "]" {
				return errors.New("Invalid term!")
			}

			(*index)++
		} else if slices.Contains([]string{"(", "."}, tokens[*index].Value) {
			if err := compileSubroutineCall(output, tokens, index, identifier); err != nil {
				return err
			}
		}
	}

	return nil
}

func compileExpression(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if err := compileTerm(output, tokens, index); err != nil {
		return err
	}

	if slices.Contains([]string{"+", "-", "*", "/", "&", "|", "<", ">", "="}, tokens[*index].Value) {
		op := tokens[*index].Value

		(*index)++

		if err := compileTerm(output, tokens, index); err != nil {
			return err
		}

		switch op {
		case "+":
			code.WriteArithmeticLogical(output, code.ADD)
		case "-":
			code.WriteArithmeticLogical(output, code.SUB)
		case "*":
			code.WriteCall(output, "Math.multiply", 2)
		case "/":
			code.WriteCall(output, "Math.divide", 2)
		case "&":
			code.WriteArithmeticLogical(output, code.AND)
		case "|":
			code.WriteArithmeticLogical(output, code.OR)
		case "<":
			code.WriteArithmeticLogical(output, code.LT)
		case ">":
			code.WriteArithmeticLogical(output, code.GT)
		case "=":
			code.WriteArithmeticLogical(output, code.EQ)
		}
	}

	return nil
}

func compileExpressionList(output *strings.Builder, tokens []tokenizer.Token, index *int) (int, error) {
	args := 0

	if slices.Contains([]tokenizer.TokenType{tokenizer.IDENTIFIER, tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) || slices.Contains([]string{"true", "false", "null", "this", "~", "-", "("}, tokens[*index].Value) {
		if err := compileExpression(output, tokens, index); err != nil {
			return 0, err
		}

		args++

		if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
			(*index)++

			more, err := compileExpressionList(output, tokens, index)

			if err != nil {
				return 0, err
			}

			args += more
		}
	}

	return args, nil
}

func compileParameterList(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) || !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return nil
	}

	variableType := tokens[*index].Value
	kind := table.ARG

	(*index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	table.AppendVariable(&subroutineSymbolTable, tokens[*index].Value, variableType, kind)
	(*index)++

	if tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		(*index)++

		if err := compileParameterList(output, tokens, index); err != nil {
			return err
		}
	}

	return nil
}

func compileVariableDeclaration(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "var" {
		return nil
	}

	(*index)++

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid variable type name!")
	}

	variableType := tokens[*index].Value

	(*index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	table.AppendVariable(&subroutineSymbolTable, tokens[*index].Value, variableType, table.VAR)
	(*index)++

	for tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		(*index)++

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid variable name!")
		}

		variableType := tokens[*index].Value

		table.AppendVariable(&subroutineSymbolTable, tokens[*index].Value, variableType, table.VAR)
		(*index)++
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	(*index)++

	return compileVariableDeclaration(output, tokens, index)
}

func compileSubroutineCall(output *strings.Builder, tokens []tokenizer.Token, index *int, identifier string) error {
	var class string
	var function string
	var subroutine string

	instanceCallingMethod := false

	switch tokens[*index].Value {
	case ".":
		variable, found := table.GetVariable([]*map[string]table.Variable{&subroutineSymbolTable, &classSymbolTable}, identifier)

		if found {
			class = variable.Type
		} else {
			class = identifier
		}

		(*index)++

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid subroutine name!")
		}

		subroutine = tokens[*index].Value
		function = class + "." + subroutine

		// Assumes that the called subroutine of a class instance is a method.
		if found {
			code.WritePush(output, segmentFromVariableKind(variable.Kind), variable.Count)
			instanceCallingMethod = true
		}

		(*index)++
	case "(":
		subroutine = identifier
		function = className + "." + identifier

		if isMethodName(tokens, subroutine) {
			code.WritePush(output, code.POINTER, 0)
		}
	}

	if tokens[*index].Value != "(" {
		return errors.New("Missing subroutine call opening parenthese!")
	}

	(*index)++

	args, err := compileExpressionList(output, tokens, index)

	if err != nil {
		return err
	}

	if tokens[*index].Value != ")" {
		return errors.New("Missing subroutine call closing parenthese!")
	}

	if isMethodName(tokens, subroutine) || instanceCallingMethod {
		args++
	}

	code.WriteCall(output, function, args)

	(*index)++

	return nil
}

func compileLetStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "let" {
		return errors.New("Invalid let statement!")
	}

	(*index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	variable, found := table.GetVariable([]*map[string]table.Variable{&subroutineSymbolTable, &classSymbolTable}, tokens[*index].Value)

	if !found {
		return errors.New("Invalid identifier!")
	}

	(*index)++

	if tokens[*index].Value == "[" {
		(*index)++

		if err := compileExpression(output, tokens, index); err != nil {
			return err
		}

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "]" {
			return errors.New("Invalid expression!")
		}

		(*index)++
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "=" {
		return errors.New("Missing assignment!")
	}

	(*index)++

	if err := compileExpression(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	code.WritePop(output, segmentFromVariableKind(variable.Kind), variable.Count)
	(*index)++

	return nil
}

func compileIfStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	endLabel := "IF_END_" + fmt.Sprint(ifLabelIndex)
	elseLabel := "IF_ELSE_" + fmt.Sprint(ifLabelIndex)

	ifLabelIndex++

	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "if" {
		return errors.New("Invalid if statement!")
	}

	(*index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing if statement opening parenthese!")
	}

	(*index)++

	if err := compileExpression(output, tokens, index); err != nil {
		return err
	}

	code.WriteArithmeticLogical(output, code.NOT)
	code.WriteIfGoto(output, elseLabel)

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing if statement closing parenthese!")
	}

	(*index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing if statement opening curly brace!")
	}

	(*index)++

	if err := compileStatements(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing if statement closing curly brace!")
	}

	(*index)++

	if tokens[*index].Value != "else" {
		code.WriteLabel(output, elseLabel)
		code.WriteGoto(output, endLabel)
	} else if tokens[*index].Type == tokenizer.KEYWORD && tokens[*index].Value == "else" {
		(*index)++

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
			return errors.New("Missing if statement opening curly brace!")
		}

		(*index)++

		code.WriteGoto(output, endLabel)
		code.WriteLabel(output, elseLabel)

		if err := compileStatements(output, tokens, index); err != nil {
			return err
		}

		if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
			return errors.New("Missing if statement closing curly brace!")
		}

		(*index)++
	}

	code.WriteLabel(output, endLabel)

	return nil
}

func compileWhileStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	startLabel := "WHILE_START_" + fmt.Sprint(whileLabelIndex)
	endLabel := "WHILE_END_" + fmt.Sprint(whileLabelIndex)

	whileLabelIndex++

	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "while" {
		return errors.New("Invalid while statement!")
	}

	(*index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing while statement opening parenthese!")
	}

	(*index)++

	code.WriteLabel(output, startLabel)

	if err := compileExpression(output, tokens, index); err != nil {
		return err
	}

	code.WriteArithmeticLogical(output, code.NOT)
	code.WriteIfGoto(output, endLabel)

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing while statement closing parenthese!")
	}

	(*index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing while statement opening curly brace!")
	}

	(*index)++

	if err := compileStatements(output, tokens, index); err != nil {
		return err
	}

	code.WriteGoto(output, startLabel)

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing while statement closing curly brace!")
	}

	code.WriteLabel(output, endLabel)

	(*index)++

	return nil
}

func compileDoStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "do" {
		return errors.New("Invalid do statement!")
	}

	(*index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	identifier := tokens[*index].Value

	(*index)++

	if err := compileSubroutineCall(output, tokens, index, identifier); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	(*index)++
	code.WritePop(output, code.TEMP, 0)

	return nil
}

func compileReturnStatement(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || tokens[*index].Value != "return" {
		return errors.New("Invalid return statement!")
	}

	(*index)++

	if slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER, tokenizer.INT_CONST, tokenizer.STR_CONST}, tokens[*index].Type) {
		if err := compileExpression(output, tokens, index); err != nil {
			return err
		}
	} else {
		code.WritePush(output, code.CONSTANT, 0)
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	code.WriteReturn(output)
	(*index)++

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

func compileSubroutineBody(output *strings.Builder, tokens []tokenizer.Token, index *int, function string, isMethod bool, isConstructor bool) error {
	if tokens[*index].Type == tokenizer.KEYWORD && tokens[*index].Value == "var" {
		if err := compileVariableDeclaration(output, tokens, index); err != nil {
			return err
		}
	}

	code.WriteFunction(output, className+"."+function, table.CountVariables(&subroutineSymbolTable, table.VAR)+1)

	if isConstructor {
		code.WritePush(output, code.CONSTANT, table.CountVariables(&classSymbolTable, table.FIELD)+1)
		code.WriteCall(output, "Memory.alloc", 1)
		code.WritePop(output, code.POINTER, 0)
	}

	if isMethod {
		code.WritePush(output, code.ARGUMENT, 0)
		code.WritePop(output, code.POINTER, 0)
	}

	if err := compileStatements(output, tokens, index); err != nil {
		return err
	}

	return nil
}

func compileSubroutineDeclaration(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	subroutineSymbolTable = make(map[string]table.Variable)

	if tokens[*index].Type != tokenizer.KEYWORD || !slices.Contains([]string{"constructor", "method", "function"}, tokens[*index].Value) {
		return nil
	}

	isMethod := tokens[*index].Value == "method"
	isConstructor := tokens[*index].Value == "constructor"

	(*index)++

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"void", "int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid subroutine return type!")
	}

	(*index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid subroutine name!")
	}

	function := tokens[*index].Value

	(*index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "(" {
		return errors.New("Missing subroutine opening parenthese!")
	}

	(*index)++

	if isMethod {
		variableType := className
		kind := table.ARG

		table.AppendVariable(&subroutineSymbolTable, "this", variableType, kind)
	}

	if err := compileParameterList(output, tokens, index); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ")" {
		return errors.New("Missing subroutine closing parenthese!")
	}

	(*index)++

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "{" {
		return errors.New("Missing subroutine opening curly brace!")
	}

	(*index)++

	if err := compileSubroutineBody(output, tokens, index, function, isMethod, isConstructor); err != nil {
		return err
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != "}" {
		return errors.New("Missing subroutine closing curly brace!")
	}

	(*index)++

	return compileSubroutineDeclaration(output, tokens, index)
}

func compileClassVarDec(output *strings.Builder, tokens []tokenizer.Token, index *int) error {
	if tokens[*index].Type != tokenizer.KEYWORD || !slices.Contains([]string{"static", "field"}, tokens[*index].Value) {
		return nil
	}

	var kind table.VariableKind

	if tokens[*index].Value == "static" {
		kind = table.STATIC
	} else {
		kind = table.FIELD
	}

	(*index)++

	if !slices.Contains([]tokenizer.TokenType{tokenizer.KEYWORD, tokenizer.IDENTIFIER}, tokens[*index].Type) && !slices.Contains([]string{"int", "char", "boolean"}, tokens[*index].Value) {
		return errors.New("Invalid variable type name!")
	}

	variableType := tokens[*index].Value

	(*index)++

	if tokens[*index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid variable name!")
	}

	table.AppendVariable(&classSymbolTable, tokens[*index].Value, variableType, kind)
	(*index)++

	for tokens[*index].Type == tokenizer.SYMBOL && tokens[*index].Value == "," {
		(*index)++

		if tokens[*index].Type != tokenizer.IDENTIFIER {
			return errors.New("Invalid variable name!")
		}

		table.AppendVariable(&classSymbolTable, tokens[*index].Value, variableType, kind)
		(*index)++
	}

	if tokens[*index].Type != tokenizer.SYMBOL || tokens[*index].Value != ";" {
		return errors.New("Missing semicolon!")
	}

	(*index)++

	return compileClassVarDec(output, tokens, index)
}

func compileClass(output *strings.Builder, tokens []tokenizer.Token) error {
	index := 0

	classSymbolTable = make(map[string]table.Variable)

	if tokens[index].Type != tokenizer.KEYWORD || tokens[index].Value != "class" {
		return errors.New("Jack file must contain one class!")
	}

	index++

	if tokens[index].Type != tokenizer.IDENTIFIER {
		return errors.New("Invalid class name!")
	}

	className = tokens[index].Value

	index++

	if tokens[index].Type != tokenizer.SYMBOL || tokens[index].Value != "{" {
		return errors.New("Missing class opening curly brace!")
	}

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

	index++

	return nil
}

func Compile(tokens []tokenizer.Token) (string, error) {
	var output strings.Builder

	if err := compileClass(&output, tokens); err != nil {
		return "", err
	}

	return output.String(), nil
}
