package executors

import (
	"log"
	"time"

	"github.com/BruhaBruh/bruhautomate/internal/flag"
)

const SleepExecutorType = "sleep"

type SleepExecutor struct {
	Type     string        `yaml:"type"`
	Duration time.Duration `yaml:"duration"`
}

var _ Executor = (*SleepExecutor)(nil)

func (e *SleepExecutor) Execute(
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
	log.Printf("[SleepExecutor] Sleep for %s\n", e.Duration.String())
	time.Sleep(e.Duration)
	log.Println("[SleepExecutor] Done")
	return nil
}

func (e *SleepExecutor) Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag {
	return []flag.Flag{}
}
