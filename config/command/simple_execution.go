package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type SimpleExecution struct {
	Type        string `yaml:"type"`
	Instruction string `yaml:"instruction"`
	Plain       bool   `yaml:"plain,omitempty"`
}

func (e *SimpleExecution) Execute(_ []Command, args []string, flags map[string]string) error {
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Command failed: %w", err)
	}

	return nil
}
