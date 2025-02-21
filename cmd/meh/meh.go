package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cpcf/meh/internal/client"
)

func main() {
	opts := parseFlags()

	// Build a final query from CLI arguments prepended to any piped input.
	finalQuery := buildQuery(flag.Args(), readStdin())
	if finalQuery != "" {
		opts.QueryArgs = []string{finalQuery}
	}

	if err := client.RunApp(opts); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

// parseFlags parses the command-line flags and returns a client.Options struct.
func parseFlags() client.Options {
	interactive := flag.Bool("i", false, "Run in interactive mode")
	filePath := flag.String("f", "", "Read input from a file")
	configFlag := flag.Bool("config", false, "Edit config settings")
	selectModel := flag.Bool("m", false, "Select a default model")
	urlFlag := flag.String("url", "", "Base URL for the LLM API")
	flag.Parse()

	return client.Options{
		Interactive: *interactive,
		FilePath:    *filePath,
		Config:      *configFlag,
		SelectModel: *selectModel,
		URL:         *urlFlag,
	}
}

// readStdin returns trimmed piped input from STDIN.
// If no piped input exists, it returns an empty string.
func readStdin() string {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Error stating stdin: %v", err)
	}
	// Check if input is piped (i.e. not coming from a terminal)
	if info.Mode()&os.ModeCharDevice != 0 {
		return ""
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading stdin: %v", err)
	}
	return strings.TrimSpace(string(data))
}

// buildQuery concatenates CLI arguments (prepended) with any STDIN input.
func buildQuery(cliArgs []string, stdinInput string) string {
	var sb strings.Builder

	if len(cliArgs) > 0 {
		sb.WriteString(strings.Join(cliArgs, " "))
	}
	if stdinInput != "" {
		if sb.Len() > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteString(stdinInput)
	}
	return sb.String()
}
