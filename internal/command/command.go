package command

import (
	"slices"

	"github.com/BruhaBruh/bruhautomate/internal/executor"
	"github.com/BruhaBruh/bruhautomate/internal/flag"
	"gopkg.in/yaml.v3"
)

type Command struct {
	Name                      string              `yaml:"-"`
	Description               string              `yaml:"description"`
	Aliases                   []string            `yaml:"aliases,omitempty"`
	Flags                     []flag.Flag         `yaml:"flags,omitempty"`
	Executors                 []executor.Executor `yaml:"executors"`
	getCommands               func() []Command    `yaml:"-"`
	*executor.ExecutorManager `yaml:"-"`
}

func (c *Command) Update(name string, getCommands func() []Command) *Command {
	return &Command{
		Name:        name,
		Description: c.Description,
		Aliases:     c.Aliases,
		Flags:       c.Flags,
		Executors:   c.Executors,
		getCommands: getCommands,
		ExecutorManager: executor.NewExecutorManager(c.Executors, func(commandName string) *executor.ExecutorManager {
			for _, command := range getCommands() {
				if command.Is(commandName) {
					return command.ExecutorManager
				}
			}
			return nil
		}),
	}
}

func (c *Command) Is(cmd string) bool {
	return c.Name == cmd || slices.Contains(c.Aliases, cmd)
}

func (c *Command) ComposeFlags() *flag.Flags {
	flags := flag.NewFlags(c.Flags...)
	for _, ex := range c.Executors {
		executorFlags := ex.Flags(func(name string) []flag.Flag {
			for _, command := range c.getCommands() {
				if command.Is(name) {
					return command.ComposeFlags().Flags()
				}
			}
			return make([]flag.Flag, 0)
		})
		flags.AddFlags(executorFlags...)
	}

	return flags
}

func (c *Command) UnmarshalYAML(value *yaml.Node) error {
	type Alias Command
	aux := struct {
		Description string      `yaml:"description"`
		Aliases     []string    `yaml:"aliases,omitempty"`
		Flags       []flag.Flag `yaml:"flags,omitempty"`
		Executors   []yaml.Node `yaml:"executors"`
	}{}

	if err := value.Decode(&aux); err != nil {
		return err
	}

	c.Description = aux.Description
	c.Aliases = aux.Aliases
	c.Flags = aux.Flags
	c.Executors = nil

	for i := range aux.Executors {
		node := &aux.Executors[i]
		if node.Kind == 0 {
			continue
		}

		exec, err := executor.UnmarshalExecutor(node)
		if err != nil {
			return err
		}
		c.Executors = append(c.Executors, exec)
	}

	return nil
}
