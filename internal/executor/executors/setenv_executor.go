package executors

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/BruhaBruh/bruhautomate/internal/errors"
	"github.com/BruhaBruh/bruhautomate/internal/flag"
	"github.com/BruhaBruh/bruhautomate/internal/instruction"
)

const SetEnvExecutorType = "setenv"

var (
	ErrFailSetEnv = errors.New("fail set env")
)

type SetEnvExecutor struct {
	Type        string `yaml:"type"`
	Name        string `yaml:"name"`
	Instruction string `yaml:"instruction"`
}

var _ Executor = (*SetEnvExecutor)(nil)

func (e *SetEnvExecutor) Execute(
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
	log.Printf("[SetEnvExecutor] Execute `%s` and set env `%s`\n", shellInstruction, e.Name)

	cmd := exec.Command("sh", "-c", shellInstruction)
	out, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, ErrFailRunCommand)
	}

	value := strings.TrimSpace(string(out))
	err = os.Setenv(e.Name, value)
	if err != nil {
		return errors.Wrap(err, ErrFailSetEnv)
	}

	return nil
}

func (e *SetEnvExecutor) Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag {
	return []flag.Flag{}
}
