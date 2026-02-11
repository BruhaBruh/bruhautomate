package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type EnvExecution struct {
	Type        string `yaml:"type"`
	Name        string `yaml:"name"`
	Instruction string `yaml:"instruction"`
	Plain       bool   `yaml:"plain,omitempty"`
}

func (e *EnvExecution) Execute(_ []Command, args []string, flags map[string]string) error {
	cmdInstruction := e.Instruction
	if !e.Plain {
		cmdInstruction = formatInsturction(e.Instruction, args, flags)
	}
	log.Printf("Execute: `%s`\n", cmdInstruction)

	parts := strings.Fields(cmdInstruction)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Command failed: %w", err)
	}

	value := strings.TrimSpace(string(out))
	err = os.Setenv(e.Name, value)
	if err != nil {
		return fmt.Errorf("Failed to set env: %w", err)
	}

	return nil
}
