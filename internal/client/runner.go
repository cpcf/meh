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
		return fmt.Errorf("error loading config: %w", err)
	}

	// Select the active API.
	activeAPI, conf, err := selectAPI(conf, opts.URL, opts.Role)
	if err != nil {
		return err
	}
	// Initialize the API client.
	api := ollama.NewAPI(activeAPI.APIURL, activeAPI.SystemPrompt)
	models := api.Models()
	if len(models) == 0 {
		return fmt.Errorf("no models found")
	}

	// Ensure a valid default model is selected.
	if activeAPI.DefaultModel == "" || !contains(models, activeAPI.DefaultModel) || opts.SelectModel {
		if err := selectModelInteractively(api, models, &activeAPI); err != nil {
			return err
		}
		updateAPIConfig(conf, activeAPI)
		if err := SaveConfig(conf); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}
		// If we were only selecting a model, exit after updating.
		if opts.SelectModel {
			return nil
		}
	} else {
		api.SelectModel(activeAPI.DefaultModel)
	}

	switch {
	case opts.Interactive:
		return runInteractive(api)
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
		return runInteractive(api)
	}
}

func selectAPI(conf *Config, urlFlag, roleFlag string) (APIConfig, *Config, error) {
	var activeAPI APIConfig
	if roleFlag != "" {
		role, found := FindRole(*conf, roleFlag)
		if !found {
			fmt.Printf("Role '%s' does not exist. Let's create it.\n", roleFlag)
			var err error
			role, err = CreateRole(roleFlag, *conf)
			if err != nil {
				return activeAPI, nil, err
			}
			conf.Roles = append(conf.Roles, role)
			if err := SaveConfig(conf); err != nil {
				return activeAPI, nil, fmt.Errorf("error saving config: %w", err)
			}
			fmt.Printf("Role %v created successfully.\n", role)
		}
		activeAPI = APIConfig{
			APIURL:       role.APIURL,
			DefaultModel: role.Model,
			SystemPrompt: role.SystemPrompt,
		}
	} else if urlFlag != "" {
		found := false
		for _, apiConf := range conf.APIs {
			if apiConf.APIURL == urlFlag {
				activeAPI = apiConf
				found = true
				break
			}
		}
		if !found {
			activeAPI = APIConfig{APIURL: urlFlag}
			conf.APIs = append(conf.APIs, activeAPI)
			if err := SaveConfig(conf); err != nil {
				return activeAPI, conf, fmt.Errorf("error saving config: %w", err)
			}
		}
	} else if len(conf.APIs) > 0 {
		activeAPI = conf.APIs[0]
	} else {
		return activeAPI, conf, fmt.Errorf("no APIs configured")
	}
	return activeAPI, conf, nil
}

// updateAPIConfig updates the default model for the active API in the config.
func updateAPIConfig(conf *Config, activeAPI APIConfig) {
	for i, a := range conf.APIs {
		if a.APIURL == activeAPI.APIURL {
			conf.APIs[i].DefaultModel = activeAPI.DefaultModel
			return
		}
	}
	// Fallback: update the first API if not found.
	if len(conf.APIs) > 0 {
		conf.APIs[0].DefaultModel = activeAPI.DefaultModel
	}
}

type API interface {
	Models() []string
	SelectModel(model string)
	Chat(query string, results chan string, flag bool)
	Prompt(query string, results chan string, flag bool)
}

func selectModelInteractively(api API, models []string, activeAPI *APIConfig) error {
	if len(models) == 1 {
		activeAPI.DefaultModel = models[0]
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
	activeAPI.DefaultModel = models[modelIndex]
	api.SelectModel(activeAPI.DefaultModel)
	return nil
}

// runInteractive runs a simple interactive chat.
func runInteractive(api API) error {
	fmt.Println("Entering interactive mode. Type your messages below.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		results := make(chan string)
		go api.Chat(line, results, true)
		for res := range results {
			fmt.Print(res)
		}
		fmt.Println()
	}
	return nil
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
