package client

import (
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

// APIConfig holds configuration for a single API.
type APIConfig struct {
	APIURL       string `yaml:"api_url"`
	DefaultModel string `yaml:"default_model"`
}

// Config holds a collection of API configurations.
type Config struct {
	APIs []APIConfig `yaml:"apis"`
}

var confpath = "/.config/.meh/config.yml"

// loadConfig reads the config file from the user's home directory.
// If no config exists, it returns a default configuration.
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := home + confpath
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConf := &Config{

			APIs: []APIConfig{
				{
					APIURL:       "http://minty.local:11434/api",
					DefaultModel: "",
				},
			},
		}
		if err := SaveConfig(defaultConf); err != nil {
			return nil, err
		}
		return defaultConf, nil
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

// saveConfig writes the configuration to the config file.
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
