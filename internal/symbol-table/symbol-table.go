package table

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

func AppendVariable(symbolTable *map[string]Variable, name string, variableType string, kind VariableKind) {
	(*symbolTable)[name] = Variable{Type: variableType, Kind: kind, Count: CountVariables(symbolTable, kind) + 1, IsDeclared: true}
}
