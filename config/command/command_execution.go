package command

type CommandExecution struct {
	Type    string   `yaml:"type"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args,omitempty"`
}

func (e *CommandExecution) Execute(commands []Command, args []string, flags map[string]string) error {
	cmd := FindCommand(commands, e.Command)
	cmdArgs := make([]string, 0, len(e.Args)+len(args))
	cmdArgs = append(cmdArgs, e.Args...)
	cmdArgs = append(cmdArgs, args...)
	return cmd.ExecuteAF(commands, cmdArgs, flags)
}
