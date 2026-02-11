package command

import (
	"slices"
	"strings"
)

type ConditionExecution struct {
	Type             string     `yaml:"type"`
	Flag             string     `yaml:"flag"`
	OnTrueExecution  Executable `yaml:"onTrueExecution"`
	OnFalseExecution Executable `yaml:"onFalseExecution"`
}

func (e *ConditionExecution) Execute(commands []Command, args []string, flags map[string]string) error {
	value, ok := flags[e.Flag]
	if !ok || !e.IsTruth(strings.ToLower(value)) {
		return e.OnFalseExecution.Execute(commands, args, flags)
	}
	return e.OnTrueExecution.Execute(commands, args, flags)
}

func (e *ConditionExecution) IsTruth(value string) bool {
	truthValues := []string{"1", "yes", "on", "true"}
	return slices.Contains(truthValues, value)
}
