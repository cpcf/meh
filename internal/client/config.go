package client

import (
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type ErrNoConfig error

type Config struct {
	DefaultPersona string    `yaml:"default_persona"`
	Personas       []Persona `yaml:"personas"`
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

func (c *Config) FindPersona(name string) (Persona, bool) {
	for _, persona := range c.Personas {
		if persona.Name == name {
			return persona, true
		}
	}
	return Persona{}, false
}

func (c *Config) AddPersona(persona Persona, setDefault bool) {
	c.Personas = append(c.Personas, persona)
	if setDefault {
		c.DefaultPersona = persona.Name
	}
	SaveConfig(c)
}

func (c *Config) LoadDefaultPersona(opts Options) (Persona, bool) {
	havePersona := false
	persona, ok := c.FindPersona(c.DefaultPersona)
	if ok {
		havePersona = true
	}
	if opts.Persona != "" {
		specifiedPersona, ok := c.FindPersona(opts.Persona)
		if ok {
			havePersona = true
			persona = specifiedPersona
		}
	}
	return persona, havePersona
}
