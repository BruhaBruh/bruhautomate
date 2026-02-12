package executor

import (
	"github.com/BruhaBruh/bruhautomate/internal/errors"
	"github.com/BruhaBruh/bruhautomate/internal/flag"
)

var (
	ErrFailExecute     = errors.New("fail execute")
	ErrCommandNotFound = errors.New("command not found")
)

type ExecutorManager struct {
	executors                    []Executor
	commandExecutorManagerGetter func(commandName string) *ExecutorManager
}

func NewExecutorManager(
	executors []Executor,
	commandExecutorManagerGetter func(commandName string) *ExecutorManager,
) *ExecutorManager {
	return &ExecutorManager{
		executors:                    executors,
		commandExecutorManagerGetter: commandExecutorManagerGetter,
	}
}

func (e *ExecutorManager) Execute(args []string, flags *flag.Flags, commandFlags *flag.Flags) error {
	for _, ex := range e.executors {
		if err := ex.Execute(args, flags, commandFlags, e.executeCommand); err != nil {
			return errors.Wrap(err, ErrFailExecute)
		}
	}
	return nil
}

func (e *ExecutorManager) executeCommand(
	commandName string,
	args []string,
	flags *flag.Flags,
	commandFlags *flag.Flags,
) error {
	commandExecutorManager := e.commandExecutorManagerGetter(commandName)
	if commandExecutorManager == nil {
		return ErrCommandNotFound
	}
	return commandExecutorManager.Execute(args, flags, commandFlags)
}
