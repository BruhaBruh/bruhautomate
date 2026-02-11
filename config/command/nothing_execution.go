package command

type NothingExecution struct {
	Type string `yaml:"type"`
}

func (e *NothingExecution) Execute(commands []Command, args []string, flags map[string]string) error {
	return nil
}
