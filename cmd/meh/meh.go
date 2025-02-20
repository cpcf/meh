package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cpcf/meh/internal/client"
	"github.com/cpcf/meh/internal/server"
)

func main() {
	serve := flag.String("s", "8080", "Server address")
	interactive := flag.Bool("i", false, "Run in interactive mode")
	filePath := flag.String("f", "", "Read input from a file")
	config := flag.Bool("config", false, "Show config settings")
	flag.Parse()

	if *serve != "" {
		http.HandleFunc("/api/llm", server.LlmHandler)

		log.Printf("Server started on port %s", *serve)
		log.Fatal(http.ListenAndServe(":"+*serve, nil))
		return
	}

	if *config {
		fmt.Println("Config settings: (Placeholder for future configuration)")
		return
	}

	if *interactive {
		client.InteractiveMode()
		return
	}

	if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			log.Fatal("Error opening file:", err)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			log.Fatal("Error reading file:", err)
		}
		client.SendRequest(string(content))
		return
	}

	if flag.NArg() > 0 {
		client.SendRequest(strings.Join(flag.Args(), " "))
		return
	}

	fmt.Println("Usage: meh [options] <query>")
	flag.PrintDefaults()
}
