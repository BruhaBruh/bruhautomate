package executors

import (
	"github.com/BruhaBruh/bruhautomate/internal/flag"
)

const CommandExecutorType = "command"

type CommandExecutor struct {
	Type    string   `yaml:"type"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args,omitempty"`
}

var _ Executor = (*CommandExecutor)(nil)

func (e *CommandExecutor) Execute(
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
	cmdArgs := make([]string, 0, len(args)+len(e.Args))
	cmdArgs = append(cmdArgs, e.Args...)
	cmdArgs = append(cmdArgs, args...)
	return executeCommand(e.Command, cmdArgs, flags, commandFlags)
}

func (e *CommandExecutor) Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag {
	return commandFlagsGetter(e.Command)
}
