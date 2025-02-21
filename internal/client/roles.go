package client

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

// CreateRole interactively prompts the user to create a new role.
func CreateRole(roleName string, conf Config) (Role, error) {
	reader := bufio.NewReader(os.Stdin)

	// List available APIs.
	fmt.Println("Select a configured API for this role:")
	for i, apiConf := range conf.APIs {
		fmt.Printf("%d: %s (default model: %s)\n", i, apiConf.APIURL, apiConf.DefaultModel)
	}
	fmt.Print("Enter the number of the API: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return Role{}, err
	}
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 0 || index >= len(conf.APIs) {
		return Role{}, fmt.Errorf("invalid API selection")
	}
	selectedAPI := conf.APIs[index]

	// Instantiate API client.
	apiInstance := ollama.NewAPI(selectedAPI.APIURL)
	models := apiInstance.Models()
	if len(models) == 0 {
		return Role{}, fmt.Errorf("no models found for the selected API")
	}

	// Let the user select a model.
	var selectedModel string
	if len(models) == 1 {
		selectedModel = models[0]
		fmt.Printf("Only one model available: %s selected.\n", selectedModel)
	} else {
		fmt.Println("Available models:")
		for i, model := range models {
			fmt.Printf("%d: %s\n", i, model)
		}
		fmt.Print("Select a model (by number): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return Role{}, err
		}
		input = strings.TrimSpace(input)
		idx, err := strconv.Atoi(input)
		if err != nil || idx < 0 || idx >= len(models) {
			return Role{}, fmt.Errorf("invalid model selection")
		}
		selectedModel = models[idx]
	}

	// Prompt for an optional system prompt.
	fmt.Print("Enter an optional system prompt (or press Enter to skip): ")
	sysPrompt, err := reader.ReadString('\n')
	if err != nil {
		return Role{}, err
	}
	sysPrompt = strings.TrimSpace(sysPrompt)

	return Role{
		Name:         roleName,
		APIURL:       selectedAPI.APIURL,
		Model:        selectedModel,
		SystemPrompt: sysPrompt,
	}, nil
}
