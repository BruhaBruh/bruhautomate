package executors

import "github.com/BruhaBruh/bruhautomate/internal/flag"

type Executor interface {
	Execute(
		args []string,
		flags *flag.Flags,
		commandFlags *flag.Flags,
		executeCommand func(
			name string,
			args []string,
			flags *flag.Flags,
			commandFlags *flag.Flags,
		) error,
	) error
	Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag
}
