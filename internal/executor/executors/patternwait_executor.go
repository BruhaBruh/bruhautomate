package executors

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/BruhaBruh/bruhautomate/internal/errors"
	"github.com/BruhaBruh/bruhautomate/internal/flag"
	"github.com/BruhaBruh/bruhautomate/internal/instruction"
)

const PatternWaitExecutorType = "patternwait"

var (
	ErrFailCreateStdoutPipe = errors.New("fail create stdout pipe")
	ErrFailStartCommand     = errors.New("fail start command")
	ErrFailScan             = errors.New("fail scan")
)

type PatternWaitExecutor struct {
	Type        string `yaml:"type"`
	Instruction string `yaml:"instruction"`
	Pattern     string `yaml:"pattern"`
}

var _ Executor = (*PatternWaitExecutor)(nil)

func (e *PatternWaitExecutor) Execute(
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
	log.Printf("[PatternWaitExecutor] Execute `%s` and wait `%s` pattern\n", shellInstruction, e.Pattern)

	parts := strings.Fields(shellInstruction)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, ErrFailCreateStdoutPipe)
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, ErrFailStartCommand)
	}

	re := regexp.MustCompile(e.Pattern)

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()

		if re.MatchString(line) {
			cmd.Process.Kill()
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, ErrFailScan)
	}

	return nil
}

func (e *PatternWaitExecutor) Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag {
	return []flag.Flag{}
}
