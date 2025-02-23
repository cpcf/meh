package client

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpcf/meh/internal/ollama"
)

// Options represents the command-line options.
type Options struct {
	FilePath  string
	Config    bool
	Persona   string
	Help      bool
	QueryArgs []string
}

// RunApp is the main entry point into the application.
func RunApp(opts Options) error {
	switch {
	case opts.Help:
		usage()
		return nil
	case opts.Config:
		return EditConfig()
	}

	conf, err := LoadConfig()
	if err != nil {
		if _, ok := err.(ErrNoConfig); ok {
			conf = &Config{}
			SaveConfig(conf)
		}
	}

	persona, havePersona := conf.LoadDefaultPersona(opts)
	if havePersona {
		api := ollama.NewAPI(persona.APIURL, persona.Model, persona.SystemPrompt)
		switch {
		case opts.FilePath != "":
			return runFile(api, opts.FilePath)
		case len(opts.QueryArgs) > 0:
			query := strings.Join(opts.QueryArgs, " ")
			return runQuery(api, query)
		}
	}

	p := tea.NewProgram(NewMainModel(conf, persona), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
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
