package config

import (
	"os"
	"path/filepath"

	"github.com/BruhaBruh/bruhautomate/internal/command"
	"github.com/BruhaBruh/bruhautomate/internal/errors"
	"gopkg.in/yaml.v3"
)

var (
	ErrCommandNotFound  = errors.New("command not found")
	ErrFailCreateConfig = errors.New("fail create configuration file")
	ErrFailSaveConfig   = errors.New("fail save configuration")
	ErrFailLoadConfig   = errors.New("fail load configuration")
	ErrPathIsEmpty      = errors.New("path is empty")
)

type Config struct {
	path     string                     `json:"-"`
	Commands map[string]command.Command `json:"commands"`
}

func New(path string) *Config {
	return &Config{
		path:     path,
		Commands: make(map[string]command.Command),
	}
}

func (c *Config) CommandsSlice() []command.Command {
	slice := make([]command.Command, 0, len(c.Commands))
	for _, cmd := range c.Commands {
		slice = append(slice, cmd)
	}
	return slice
}

func (c *Config) Command(name string) (*command.Command, error) {
	for _, command := range c.Commands {
		if command.Is(name) {
			return &command, nil
		}
	}
	return nil, errors.Swrap(ErrCommandNotFound, name)
}

func (c *Config) LoadOrCreate() error {
	if err := c.createFileIfNotExists(); err != nil {
		return err
	}
	return c.Load()
}

func (c *Config) Load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return errors.Wrap(err, ErrFailLoadConfig)
	}
	var commands map[string]command.Command
	if err = yaml.Unmarshal(data, &commands); err != nil {
		return errors.Wrap(err, ErrFailLoadConfig)
	}
	c.Commands = commands

	for name, cmd := range c.Commands {
		c.Commands[name] = *cmd.Update(name, func() []command.Command {
			commands := make([]command.Command, 0, len(c.Commands))
			for _, cmd := range c.Commands {
				commands = append(commands, cmd)
			}
			return commands
		})
	}

	return nil
}

func (c *Config) Save() error {
	if err := c.createFileIfNotExists(); err != nil {
		return errors.Wrap(err, ErrFailSaveConfig)
	}
	data, err := yaml.Marshal(c.Commands)
	if err != nil {
		return errors.Wrap(err, ErrFailSaveConfig)
	}
	if err := os.WriteFile(c.path, data, 0644); err != nil {
		return errors.Wrap(err, ErrFailSaveConfig)
	}
	return nil
}

func (c *Config) createFileIfNotExists() error {
	if err := c.isExists(); err != nil {
		if err := c.isExistsInWorkingDirectory(); err != nil {
			if err := c.isExistsInHomeDirectory(); err != nil {
				return nil
			}
		}
	}
	if len(c.path) == 0 {
		return ErrPathIsEmpty
	}
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		dir := filepath.Dir(c.path)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.Wrap(err, errors.New("fail create directories"))
		}
		file, err := os.Create(c.path)
		if err != nil {
			return errors.Wrap(err, errors.New("fail create file"))
		}
		defer file.Close()
	}
	return nil
}

func (c *Config) isExists() error {
	if len(c.path) == 0 {
		return ErrPathIsEmpty
	}
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		return err
	}
	return nil
}

func (c *Config) isExistsInWorkingDirectory() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, ".bam.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	c.path = path
	return nil
}

func (c *Config) isExistsInHomeDirectory() error {
	hd, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(hd, ".bam.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	c.path = path
	return nil
}
