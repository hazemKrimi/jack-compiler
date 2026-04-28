package code

import "fmt"

import "strings"

type Segment int
type ArithmeticLogicalCommand int

const (
	CONSTANT Segment = iota
	ARGUMENT
	LOCAL
	STATIC
	THIS
	THAT
	POINTER
	TEMP
)

const (
	ADD ArithmeticLogicalCommand = iota
	SUB
	NEG
	EQ
	GT
	LT
	AND
	OR
	NOT
)

func SegmentText(segment Segment) string {
	switch segment {
	case CONSTANT:
		return "constant"
	case ARGUMENT:
		return "argument"
	case LOCAL:
		return "local"
	case STATIC:
		return "static"
	case THIS:
		return "this"
	case THAT:
		return "that"
	case POINTER:
		return "pointer"
	case TEMP:
		return "temp"
	}

	return "constant"
}

func WritePush(output *strings.Builder, segment Segment, index int) error {
	if _, err := output.WriteString("push " + SegmentText(segment) + " " + fmt.Sprint(index)); err != nil {
		return err
	}

	return nil
}

func WritePop(output *strings.Builder, segment Segment, index int) error {
	if _, err := output.WriteString("pop " + SegmentText(segment) + " " + fmt.Sprint(index)); err != nil {
		return err
	}

	return nil
}

func WriteArithmeticLogical(output *strings.Builder, command ArithmeticLogicalCommand) error {
	var commandText string

	switch command {
	case ADD:
		commandText = "add"
	case SUB:
		commandText = "sub"
	case NEG:
		commandText = "neg"
	case EQ:
		commandText = "eq"
	case GT:
		commandText = "gt"
	case LT:
		commandText = "lt"
	case AND:
		commandText = "and"
	case OR:
		commandText = "or"
	case NOT:
		commandText = "not"
	}

	if _, err := output.WriteString(commandText); err != nil {
		return err
	}

	return nil
}

func WriteLabel(output *strings.Builder, label string) error {
	if _, err := output.WriteString("label " + label); err != nil {
		return err
	}

	return nil
}

func WriteGoto(output *strings.Builder, label string) error {
	if _, err := output.WriteString("goto " + label); err != nil {
		return err
	}

	return nil
}

func WriteIfGoto(output *strings.Builder, label string) error {
	if _, err := output.WriteString("if-goto " + label); err != nil {
		return err
	}

	return nil
}

func WriteCall(output *strings.Builder, function string, args int) error {
	if _, err := output.WriteString("call " + function + " " + fmt.Sprint(args)); err != nil {
		return err
	}

	return nil
}

func WriteFunction(output *strings.Builder, function string, params int) error {
	if _, err := output.WriteString("function " + function + " " + fmt.Sprint(params)); err != nil {
		return err
	}

	return nil
}

func WriteReturn(output *strings.Builder) error {
	if _, err := output.WriteString("return"); err != nil {
		return err
	}

	return nil
}
