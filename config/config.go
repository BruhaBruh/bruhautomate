package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BruhaBruh/bruhautomate/config/command"
	"gopkg.in/yaml.v3"
)

type Config struct {
	path     string            `json:"-"`
	Commands []command.Command `json:"commands"`
}

func LoadOrCreate(path string) (*Config, error) {
	if len(path) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.New("Failed to get user home directory")
		}
		path = filepath.Join(homeDir, ".bam.yaml")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := &Config{
			path:     path,
			Commands: make([]command.Command, 0),
		}
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("Failed to read config")
	}
	var commands []command.Command
	err = yaml.Unmarshal(data, &commands)
	if err != nil {
		return nil, fmt.Errorf("Failed to read yaml of config: %w", err)
	}
	config := &Config{
		path:     path,
		Commands: commands,
	}
	return config, nil
}

func (c *Config) Execute(cmd string, args []string) error {
	command := command.FindCommand(c.Commands, cmd)
	if command == nil {
		return fmt.Errorf("Command not found: %s", cmd)
	}

	return command.Execute(c.Commands, args)
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c.Commands)
	if err != nil {
		return errors.New("Failed to marshal config")
	}
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		dir := filepath.Dir(c.path)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.New("Failed to create directory contains config")
		}
		file, err := os.Create(c.path)
		if err != nil {
			return errors.New("Failed to create file")
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			return errors.New("Failed to write config")
		}

		return nil
	}
	file, err := os.Open(c.path)
	if err != nil {
		return errors.New("Failed to open config")
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return errors.New("Failed to write config")
	}
	return nil
}
