package main

import (
	"flag"
	"fmt"

	"github.com/cpcf/meh/internal/ollama"
)

func main() {
	// Define a flag for the API URL with a default value.
	apiURL := flag.String("url", "http://localhost:11434/api", "Base URL for the Ollama API")
	flag.Parse()

	// Instantiate and initialize our API instance.
	var api ollama.OllamaAPI
	api.Init(*apiURL)

	fmt.Println("API initialized with URL:", *apiURL)

	// List available models.
	models := api.Models()
	if len(models) == 0 {
		fmt.Println("No models found, using default model.")
		api.SelectModel("default")
	} else {
		fmt.Println("Available models:", models)
		// Select the first available model.
		api.SelectModel(models[0])
	}

	// ----------------------------
	// Demonstrate non-streaming Chat.
	chatResults := make(chan string)
	go api.Chat("Hello, how are you?", chatResults, false)
	fmt.Println("\n[Non-streaming Chat]:")
	fmt.Println("> Hello, how are you?")
	chatResp := <-chatResults
	fmt.Println(chatResp)

	// ----------------------------
	// Demonstrate streaming Chat.
	streamChatResults := make(chan string)
	go api.Chat("What's the weather like today?", streamChatResults, true)
	fmt.Println("\n[Streaming Chat]:")
	fmt.Println("> What's the weather like today?")
	for res := range streamChatResults {
		fmt.Print(res) // multiple chunks printed as they stream in
	}
	fmt.Println() // newline after streaming chat

	// ----------------------------
	// Demonstrate non-streaming Prompt.
	promptResults := make(chan string)
	go api.Prompt("Tell me a joke.", promptResults, false)
	fmt.Println("\n[Non-streaming Prompt]:")
	fmt.Println("> Tell me a joke.")
	promptResp := <-promptResults
	fmt.Println(promptResp)

	// ----------------------------
	// Demonstrate streaming Prompt.
	streamPromptResults := make(chan string)
	go api.Prompt("Tell me a joke.", streamPromptResults, true)
	fmt.Println("\n[Streaming Prompt]:")
	fmt.Println("> Tell me a joke.")
	for res := range streamPromptResults {
		fmt.Print(res)
	}
	fmt.Println()

}
