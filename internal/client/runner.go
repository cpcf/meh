package client

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/cpcf/meh/internal/ollama"
)

// Options represents the command-line options.
type Options struct {
	Interactive bool
	FilePath    string
	Config      bool
	SelectModel bool
	URL         string
	Role        string
	Help        bool
	QueryArgs   []string
}

// RunApp is the main entry point into the application.
func RunApp(opts Options) error {
	// If the user wants to edit the config, do that and exit.
	if opts.Config {
		EditConfig()
		return nil
	}

	// Load configuration.
	conf, err := LoadConfig()
	if err != nil {
		if _, ok := err.(ErrNoConfig); ok {
			conf = &Config{}
			// If no config exists, prompt the user to create a role and then save it
			// Prompt for name
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("No config file found. Let's create our first role!\n")
			fmt.Print("Enter role name: ")
			roleName, err := reader.ReadString('\n')
			role, err := CreateRole(strings.TrimSpace(roleName), conf)
			if err != nil {
				return err
			}
			api := ollama.NewAPI(role.APIURL, role.SystemPrompt)
			api.SelectModel(role.Model)
			runInteractive(api)
			return nil
		}
		return fmt.Errorf("error loading config: %w", err)
	}

	// Select the active Role.
	role, err := selectRole(conf, opts.Role)
	if err != nil {
		return err
	}
	// Initialize the API client.
	api := ollama.NewAPI(role.APIURL, role.SystemPrompt)
	models := api.Models()
	if len(models) == 0 {
		return fmt.Errorf("no models found")
	}

	// Ensure a valid default model is selected.
	if role.Model == "" || !contains(models, role.Model) || opts.SelectModel {
		if err := selectModelInteractively(api, models, &role); err != nil {
			return err
		}
		updateRole(conf, role)
		if err := SaveConfig(conf); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}
		// If we were only updating a default model, exit after updating.
		if opts.SelectModel {
			return nil
		}
	} else {
		api.SelectModel(role.Model)
	}

	switch {
	case opts.Interactive:
		runInteractive(api)
	case opts.FilePath != "":
		return runFile(api, opts.FilePath)
	case len(opts.QueryArgs) > 0:
		query := strings.Join(opts.QueryArgs, " ")
		return runQuery(api, query)
	case opts.Help:
		usage()
		return nil
	default:
		// What should the defaut behaviour be?
		runInteractive(api)
	}
	return nil
}

func selectRole(conf *Config, roleFlag string) (Role, error) {
	if roleFlag != "" {
		role, found := FindRole(*conf, roleFlag)
		if !found {
			fmt.Printf("Role '%s' does not exist. Let's create it.\n", roleFlag)
			var err error
			role, err = CreateRole(roleFlag, conf)
			if err != nil {
				return role, err
			}
			fmt.Printf("Role %v created successfully.\n", role)
		}
		return role, nil
	} else if conf.DefaultRole != "" {
		return selectRole(conf, conf.DefaultRole)
	} else {
		return Role{}, fmt.Errorf("no roles configured")
	}
}

func updateRole(conf *Config, role Role) {
	for i, a := range conf.Roles {
		if a.Name == role.Name {
			conf.Roles[i].Model = role.Model
			return
		}
	}
}

type API interface {
	Models() []string
	SelectModel(model string)
	Chat(query string, results chan string, flag bool)
	Prompt(query string, results chan string, flag bool)
}

func selectModelInteractively(api API, models []string, role *Role) error {
	if len(models) == 1 {
		role.Model = models[0]
		api.SelectModel(models[0])
		return nil
	}
	fmt.Println("Available models:")
	for i, model := range models {
		fmt.Printf("%d: %s\n", i, model)
	}
	fmt.Print("Select a model (by number): ")
	reader := bufio.NewReader(os.Stdin)
	modelIndexStr, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	modelIndexStr = strings.TrimSpace(modelIndexStr)
	modelIndex, err := strconv.Atoi(modelIndexStr)
	if err != nil {
		return err
	}
	if modelIndex < 0 || modelIndex >= len(models) {
		return fmt.Errorf("invalid model index")
	}
	role.Model = models[modelIndex]
	api.SelectModel(role.Model)
	return nil
}

// runInteractive runs a simple interactive chat.
func runInteractive(api API) {
	ChatTui(api)
}

// runFile reads input from a file and sends it as a prompt.
func runFile(api API, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	query := string(content)
	results := make(chan string)
	go api.Prompt(query, results, true)
	for res := range results {
		fmt.Print(res)
	}
	fmt.Println()
	return nil
}

// runQuery sends a command-line query to the API.
func runQuery(api API, query string) error {
	results := make(chan string)
	go api.Prompt(query, results, true)
	for res := range results {
		fmt.Print(res)
	}
	fmt.Println()
	return nil
}

func usage() {
	fmt.Println("Usage: [options] <query>")
	flag.PrintDefaults()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
