package flag

type Flag struct {
	Name        string `yaml:"name"`
	Shortcut    string `yaml:"shortcut,omitempty"`
	Description string `yaml:"description,omitempty"`
	HasValue    bool   `yaml:"hasValue,omitempty"`
	Value       string `yaml:"-"`
}

func (f *Flag) clone(value ...string) Flag {
	flag := Flag{
		Name:        f.Name,
		Shortcut:    f.Shortcut,
		Description: f.Description,
		HasValue:    f.HasValue,
		Value:       f.Value,
	}
	if len(value) > 0 {
		flag.Value = value[0]
	}
	return flag
}
