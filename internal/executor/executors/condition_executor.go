package executors

import (
	"github.com/BruhaBruh/bruhautomate/internal/flag"
)

const ConditionExecutorType = "condition"

type ConditionExecutor struct {
	Type            string   `yaml:"type"`
	Flag            string   `yaml:"flag"`
	OnTrueExecutor  Executor `yaml:"onTrueExecutor"`
	OnFalseExecutor Executor `yaml:"onFalseExecutor"`
}

var _ Executor = (*ConditionExecutor)(nil)

func (e *ConditionExecutor) Execute(
	args []string,
	flags *flag.Flags,
	commandFlags *flag.Flags,
	executeCommand func(
		name string,
		args []string,
		flags *flag.Flags,
		commandFlags *flag.Flags,
	) error,
) error {
	if flags.HasFlag(e.Flag) {
		if e.OnTrueExecutor != nil {
			return e.OnTrueExecutor.Execute(args, flags, commandFlags, executeCommand)
		}
	} else {
		if e.OnFalseExecutor != nil {
			return e.OnFalseExecutor.Execute(args, flags, commandFlags, executeCommand)
		}
	}
	return nil
}

func (e *ConditionExecutor) Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag {
	flags := make([]flag.Flag, 0, 16)
	flags = append(flags, flag.Flag{
		Name:        e.Flag,
		Description: "Uses in condition.",
		HasValue:    false,
	})
	if e.OnTrueExecutor != nil {
		flags = append(flags, e.OnTrueExecutor.Flags(commandFlagsGetter)...)
	}
	if e.OnFalseExecutor != nil {
		flags = append(flags, e.OnFalseExecutor.Flags(commandFlagsGetter)...)
	}
	return flags
}
