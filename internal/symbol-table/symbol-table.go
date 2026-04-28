package table

import (
	"fmt"
	"strconv"
	"strings"
)

type VariableKind int

const (
	STATIC VariableKind = iota
	FIELD
	ARG
	VAR
)

type Variable struct {
	Type       string
	Kind       VariableKind
	Count      int
	IsDeclared bool
	IsUsed     bool
}

func CountVariables(symbolTable *map[string]Variable, kind VariableKind) int {
	count := -1

	for _, variable := range *symbolTable {
		if variable.Kind == kind {
			count++
		}
	}

	return count
}

func GetVariable(symbolTables []*map[string]Variable, name string) (Variable, bool) {
	for _, table := range symbolTables {
		for key, variable := range *table {
			if key == name {
				return variable, true
			}
		}
	}

	return Variable{}, false
}

func UseVariable(symbolTables []*map[string]Variable, name string) {
	for _, table := range symbolTables {
		for key, variable := range *table {
			if key == name {
				variable.IsUsed = true
				(*table)[key] = variable
				return
			}
		}
	}
}

func WriteImplicitThis(output *strings.Builder, symbolTables []*map[string]Variable) error {
	variable, found := GetVariable(symbolTables, "this")

	if found {
		tokenDefinition := "<implicitVariable> "
		tokenDefinition += "name: this, "
		tokenDefinition += "type: " + variable.Type + ", "
		tokenDefinition += "kind: " + fmt.Sprint(variable.Kind) + ", "
		tokenDefinition += "count: " + fmt.Sprint(variable.Count) + ", "
		tokenDefinition += "declared: " + strconv.FormatBool(variable.IsDeclared) + ", "
		tokenDefinition += "used: " + strconv.FormatBool(variable.IsUsed)
		tokenDefinition += "</variable>\n"

		if _, err := output.WriteString(tokenDefinition); err != nil {
			return err
		}
	}

	return nil
}

func AppendVariable(symbolTable *map[string]Variable, name string, variableType string, kind VariableKind) {
	(*symbolTable)[name] = Variable{Type: variableType, Kind: kind, Count: CountVariables(symbolTable, kind) + 1, IsDeclared: true}
}
