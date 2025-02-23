package client

import (
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type ErrNoConfig error

// Role represents a user-defined role that overrides API/model settings and may include a system prompt.
type Role struct {
	Name         string `yaml:"name"`
	APIURL       string `yaml:"api_url"`
	Model        string `yaml:"model"`
	SystemPrompt string `yaml:"system_prompt,omitempty"`
}

type Config struct {
	DefaultRole string `yaml:"default_role"`
	Roles       []Role `yaml:"roles"`
}

var confpath = "/.config/.meh/config.yml"

// loadConfig reads the config file from the user's home directory.
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := home + confpath
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, ErrNoConfig(err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func SaveConfig(conf *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := home + confpath
	err = os.MkdirAll(home+"/.config/.meh", os.ModePerm)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func EditConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := home + confpath
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Config) FindRole(name string) (Role, bool) {
	for _, role := range c.Roles {
		if role.Name == name {
			return role, true
		}
	}
	return Role{}, false
}
