package flag

import (
	"github.com/BruhaBruh/bruhautomate/internal/errors"
)

var ErrFlagNotFound = errors.New("flag not found")

type Flags struct {
	flags map[string]Flag
}

func NewFlags(current ...Flag) *Flags {
	flags := &Flags{
		flags: make(map[string]Flag),
	}

	for _, flag := range current {
		cloned := flag.clone()
		flags.AddFlag(&cloned)
	}

	return flags
}

func (f *Flags) Flags() []Flag {
	flags := make([]Flag, 0, len(f.flags))
	for _, flag := range f.flags {
		flags = append(flags, flag)
	}
	return flags
}

func (f *Flags) AddFlag(flag *Flag, value ...string) {
	f.flags[flag.Name] = flag.clone(value...)
}

func (f *Flags) AddFlags(flags ...Flag) {
	for _, flag := range flags {
		if !f.HasFlag(flag.Name) {
			f.flags[flag.Name] = flag
		}
	}
}

func (f *Flags) Find(nameOrShortcut string) (*Flag, error) {
	for _, flag := range f.flags {
		if flag.Name == nameOrShortcut || flag.Shortcut == nameOrShortcut {
			return &flag, nil
		}
	}
	return nil, ErrFlagNotFound
}

func (f *Flags) HasFlag(nameOrShortcut string) bool {
	_, err := f.Find(nameOrShortcut)
	return err == nil
}
