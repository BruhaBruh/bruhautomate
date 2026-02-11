package command

import (
	"time"

	"github.com/BruhaBruh/bruhautomate/pkg/duration"
)

type WaitExecution struct {
	Type     string            `yaml:"type"`
	Duration duration.Duration `yaml:"duration"`
}

func (e *WaitExecution) Execute(_ []Command, _ []string, _ map[string]string) error {
	time.Sleep(e.Duration.Duration)
	return nil
}
