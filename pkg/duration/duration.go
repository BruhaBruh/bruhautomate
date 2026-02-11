package duration

import (
	"errors"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalYAML() (any, error) {
	return d.String(), nil
}

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		// Попробуем сначала как строку
		var s string
		if err := value.Decode(&s); err == nil {
			dur, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			d.Duration = dur
			return nil
		}

		var f float64
		if err := value.Decode(&f); err == nil {
			d.Duration = time.Duration(f)
			return nil
		}

		return errors.New("invalid duration value")
	default:
		return errors.New("invalid node kind for duration")
	}
}
