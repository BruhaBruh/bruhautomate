package executors

import (
	"log"
	"os"
	"os/exec"

	"github.com/BruhaBruh/bruhautomate/internal/errors"
	"github.com/BruhaBruh/bruhautomate/internal/flag"
	"github.com/BruhaBruh/bruhautomate/internal/instruction"
)

const ShellExecutorType = "shell"

var (
	ErrFailRunCommand = errors.New("fail run command")
)

type ShellExecutor struct {
	Type        string `yaml:"type"`
	Instruction string `yaml:"instruction"`
}

var _ Executor = (*ShellExecutor)(nil)

func (e *ShellExecutor) Execute(
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
	shellInstruction := instruction.Build(e.Instruction, args, flags, commandFlags)
	log.Printf("[ShellExecutor] Execute `%s`\n", shellInstruction)

	cmd := exec.Command("sh", "-c", shellInstruction)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, ErrFailRunCommand)
	}

	return nil
}

func (e *ShellExecutor) Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag {
	return []flag.Flag{}
}
