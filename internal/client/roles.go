package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cpcf/meh/internal/ollama"
)

// FindRole searches for a role by name.
func FindRole(conf Config, roleName string) (Role, bool) {
	for _, r := range conf.Roles {
		if r.Name == roleName {
			return r, true
		}
	}
	return Role{}, false
}

func CreateRole(roleName string, conf *Config) (Role, error) {
	reader := bufio.NewReader(os.Stdin)

	newRole := Role{
		Name: roleName,
	}

	// Prompt for API URL.
	fmt.Print("Enter API URL: ")
	apiURL, err := reader.ReadString('\n')
	if err != nil {
		return Role{}, err
	}
	newRole.APIURL = strings.TrimSuffix(strings.TrimSpace(apiURL), "\n")
	// Prompt for Model.
	fmt.Print("Enter model name: ")
	api := ollama.NewAPI(newRole.APIURL, "")
	err = selectModelInteractively(api, api.Models(), &newRole)
	if err != nil {
		return Role{}, err
	}

	// Prompt for an optional system prompt.
	fmt.Print("Enter an optional system prompt (or press Enter to skip): ")
	sysPrompt, err := reader.ReadString('\n')
	if err != nil {
		return Role{}, err
	}
	newRole.SystemPrompt = strings.TrimSpace(sysPrompt)

	// Add the new role to the configuration.
	conf.Roles = append(conf.Roles, newRole)

	// Prompt for a default role.
	fmt.Print("Set as default role? (y/n): ")
	defaultRole, err := reader.ReadString('\n')
	if err != nil {
		return Role{}, err
	}
	defaultRole = strings.TrimSpace(defaultRole)
	if defaultRole == "y" || defaultRole == "Y" {
		conf.DefaultRole = roleName
	}

	if err := SaveConfig(conf); err != nil {
		return Role{}, err
	}

	return newRole, nil
}
