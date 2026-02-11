package command

import (
	"fmt"

	"github.com/BruhaBruh/bruhautomate/pkg/duration"
	"gopkg.in/yaml.v3"
)

const (
	simpleExecutableType    = "simple"
	commandExecutableType   = "command"
	countdownExecutableType = "countdown"
	conditionExecutableType = "condition"
	nothingExecutableType   = "nothing"
	waitlineExecutableType  = "waitline"
	waitExecutableType      = "wait"
	envExecutableType       = "env"
)

type Executable interface {
	Execute(commands []Command, args []string, flags map[string]string) error
}

func unmarshalExecutable(node *yaml.Node) (Executable, error) {
	data, err := yaml.Marshal(node)
	if err != nil {
		return nil, err
	}

	var typeCheck struct {
		Type string `yaml:"type"`
	}
	if err := yaml.Unmarshal(data, &typeCheck); err != nil {
		return nil, err
	}

	if typeCheck.Type == "" {
		return nil, fmt.Errorf("execution type is empty")
	}

	switch typeCheck.Type {
	case simpleExecutableType:
		var ex SimpleExecution
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	case commandExecutableType:
		var ex CommandExecution
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	case countdownExecutableType:
		var alias struct {
			Type      string            `yaml:"type"`
			Time      duration.Duration `yaml:"time"`
			Execution map[string]any    `yaml:"execution"`
		}
		if err := node.Decode(&alias); err != nil {
			return nil, err
		}

		var execution Executable
		if len(alias.Execution) > 0 {
			// Конвертируем обратно в YAML и парсим
			execData, err := yaml.Marshal(alias.Execution)
			if err != nil {
				return nil, err
			}
			var execNode yaml.Node
			if err := yaml.Unmarshal(execData, &execNode); err != nil {
				return nil, err
			}
			execution, err = unmarshalExecutable(&execNode)
			if err != nil {
				return nil, err
			}
		}

		return &CountdownExecution{
			Type:      alias.Type,
			Time:      alias.Time,
			Execution: execution,
		}, nil
	case conditionExecutableType:
		var alias struct {
			Type             string         `yaml:"type"`
			Flag             string         `yaml:"flag"`
			OnTrueExecution  map[string]any `yaml:"onTrueExecution"`
			OnFalseExecution map[string]any `yaml:"onFalseExecution"`
		}
		if err := node.Decode(&alias); err != nil {
			return nil, err
		}

		var onTrue Executable
		if len(alias.OnTrueExecution) > 0 {
			execData, err := yaml.Marshal(alias.OnTrueExecution)
			if err != nil {
				return nil, err
			}
			var execNode yaml.Node
			if err := yaml.Unmarshal(execData, &execNode); err != nil {
				return nil, err
			}
			onTrue, err = unmarshalExecutable(&execNode)
			if err != nil {
				return nil, err
			}
		}

		var onFalse Executable
		if len(alias.OnFalseExecution) > 0 {
			execData, err := yaml.Marshal(alias.OnFalseExecution)
			if err != nil {
				return nil, err
			}
			var execNode yaml.Node
			if err := yaml.Unmarshal(execData, &execNode); err != nil {
				return nil, err
			}
			onFalse, err = unmarshalExecutable(&execNode)
			if err != nil {
				return nil, err
			}
		}

		return &ConditionExecution{
			Type:             alias.Type,
			Flag:             alias.Flag,
			OnTrueExecution:  onTrue,
			OnFalseExecution: onFalse,
		}, nil
	case nothingExecutableType:
		return &NothingExecution{Type: nothingExecutableType}, nil
	case waitlineExecutableType:
		var ex WaitlineExecution
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	case waitExecutableType:
		var ex WaitExecution
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	case envExecutableType:
		var ex EnvExecution
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	default:
		return nil, fmt.Errorf("unknown execution type: %s", typeCheck.Type)
	}
}
