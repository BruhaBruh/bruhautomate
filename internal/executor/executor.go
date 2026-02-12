package executor

import (
	"time"

	"github.com/BruhaBruh/bruhautomate/internal/errors"
	"github.com/BruhaBruh/bruhautomate/internal/executor/executors"
	"github.com/BruhaBruh/bruhautomate/internal/flag"
	"gopkg.in/yaml.v3"
)

var (
	ErrFailUnmarshalExecutor         = errors.New("fail unmarshal executor")
	ErrInvalidExecutorType           = errors.New("invalid executor type")
	ErrTimerExecutorChildIsUndefined = errors.New("timer executor child is undefined")
)

type Executor interface {
	Execute(
		args []string,
		flags *flag.Flags,
		commandFlags *flag.Flags,
		executeCommand func(
			name string,
			args []string,
			flags *flag.Flags,
			commandFlags *flag.Flags,
		) error,
	) error
	Flags(commandFlagsGetter func(name string) []flag.Flag) []flag.Flag
}

func UnmarshalExecutor(node *yaml.Node) (Executor, error) {
	data, err := yaml.Marshal(node)
	if err != nil {
		return nil, errors.Wrap(err, ErrFailUnmarshalExecutor)
	}

	var typeCheck struct {
		Type string `yaml:"type"`
	}
	if err := yaml.Unmarshal(data, &typeCheck); err != nil {
		return nil, errors.Wrap(err, ErrFailUnmarshalExecutor)
	}

	if typeCheck.Type == "" {
		return nil, ErrInvalidExecutorType
	}

	switch typeCheck.Type {
	case executors.CommandExecutorType:
		var ex executors.CommandExecutor
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	case executors.ConditionExecutorType:
		var alias struct {
			Type            string         `yaml:"type"`
			Flag            string         `yaml:"flag"`
			OnTrueExecutor  map[string]any `yaml:"onTrueExecutor"`
			OnFalseExecutor map[string]any `yaml:"onFalseExecutor"`
		}
		if err := node.Decode(&alias); err != nil {
			return nil, err
		}

		var onTrueExecutor Executor
		if len(alias.OnTrueExecutor) > 0 {
			execData, err := yaml.Marshal(alias.OnTrueExecutor)
			if err != nil {
				return nil, err
			}
			var execNode yaml.Node
			if err := yaml.Unmarshal(execData, &execNode); err != nil {
				return nil, err
			}
			onTrueExecutor, err = UnmarshalExecutor(&execNode)
			if err != nil {
				return nil, err
			}
		}

		var onFalseExecutor Executor
		if len(alias.OnFalseExecutor) > 0 {
			execData, err := yaml.Marshal(alias.OnFalseExecutor)
			if err != nil {
				return nil, err
			}
			var execNode yaml.Node
			if err := yaml.Unmarshal(execData, &execNode); err != nil {
				return nil, err
			}
			onFalseExecutor, err = UnmarshalExecutor(&execNode)
			if err != nil {
				return nil, err
			}
		}

		return &executors.ConditionExecutor{
			Type:            alias.Type,
			Flag:            alias.Flag,
			OnTrueExecutor:  onTrueExecutor,
			OnFalseExecutor: onFalseExecutor,
		}, nil
	case executors.PatternWaitExecutorType:
		var ex executors.PatternWaitExecutor
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	case executors.SetEnvExecutorType:
		var ex executors.SetEnvExecutor
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	case executors.ShellExecutorType:
		var ex executors.ShellExecutor
		if err := node.Decode(&ex); err != nil {
			return nil, err
		}
		return &ex, nil
	case executors.SleepExecutorType:
		var alias struct {
			Type     string `yaml:"type"`
			Duration string `yaml:"duration"`
		}
		if err := node.Decode(&alias); err != nil {
			return nil, err
		}

		duration, err := time.ParseDuration(alias.Duration)
		if err != nil {
			return nil, err
		}

		return &executors.SleepExecutor{
			Type:     alias.Type,
			Duration: duration,
		}, nil
	case executors.TimerExecutorType:
		var alias struct {
			Type     string         `yaml:"type"`
			Time     string         `yaml:"time"`
			Executor map[string]any `yaml:"executor"`
		}
		if err := node.Decode(&alias); err != nil {
			return nil, err
		}

		time, err := time.ParseDuration(alias.Time)
		if err != nil {
			return nil, err
		}

		if len(alias.Executor) == 0 {
			return nil, ErrTimerExecutorChildIsUndefined
		}
		var executor Executor
		execData, err := yaml.Marshal(alias.Executor)
		if err != nil {
			return nil, err
		}
		var execNode yaml.Node
		if err := yaml.Unmarshal(execData, &execNode); err != nil {
			return nil, err
		}
		executor, err = UnmarshalExecutor(&execNode)
		if err != nil {
			return nil, err
		}

		return &executors.TimerExecutor{
			Type:     alias.Type,
			Time:     time,
			Executor: executor,
		}, nil
	default:
		return nil, errors.Swrap(ErrInvalidExecutorType, typeCheck.Type)
	}
}
