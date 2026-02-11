package command

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type WaitlineExecution struct {
	Type        string `yaml:"type"`
	Instruction string `yaml:"instruction"`
	Plain       bool   `yaml:"plain,omitempty"`
	Pattern     string `yaml:"pattern"`
}

func (e *WaitlineExecution) Execute(commands []Command, args []string, flags map[string]string) error {
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

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
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
		log.Fatalf("Scanner has error: %v", err)
	}

	return nil
}
