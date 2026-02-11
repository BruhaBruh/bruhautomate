package command

import (
	"slices"

	"gopkg.in/yaml.v3"
)

type Command struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Aliases     []string     `yaml:"aliases,omitempty"`
	Executions  []Executable `yaml:"executions"`
}

func (c *Command) Is(cmd string) bool {
	return c.Name == cmd || slices.Contains(c.Aliases, cmd)
}

func (c *Command) Execute(commands []Command, input []string) error {
	args, flags := parseArgsAndFlags(input)
	return c.ExecuteAF(commands, args, flags)
}

func (c *Command) ExecuteAF(commands []Command, args []string, flags map[string]string) error {
	for _, e := range c.Executions {
		if err := e.Execute(commands, args, flags); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) UnmarshalYAML(value *yaml.Node) error {
	type Alias Command
	aux := struct {
		Name        string      `yaml:"name"`
		Description string      `yaml:"description"`
		Aliases     []string    `yaml:"aliases,omitempty"`
		Executions  []yaml.Node `yaml:"executions"`
	}{}

	if err := value.Decode(&aux); err != nil {
		return err
	}

	c.Name = aux.Name
	c.Description = aux.Description
	c.Aliases = aux.Aliases
	c.Executions = nil

	for i := range aux.Executions {
		node := &aux.Executions[i]
		if node.Kind == 0 {
			continue
		}

		exec, err := unmarshalExecutable(node)
		if err != nil {
			return err
		}
		c.Executions = append(c.Executions, exec)
	}

	return nil
}
