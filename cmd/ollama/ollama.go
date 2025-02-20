package main

import (
	"fmt"
	"time"

	"github.com/cpcf/meh/internal/ollama"
)

func createMessageRequest(system string, message string, model string) ollama.Request {
	sysmsg := ollama.Message{
		Role:    "system",
		Content: system,
	}
	msg := ollama.Message{
		Role:    "user",
		Content: message,
	}
	return ollama.Request{
		Model:    model,
		Stream:   false,
		Messages: []ollama.Message{sysmsg, msg},
	}
}

func createPromptRequest(prompt string, model string, stream bool) ollama.Request {
	return ollama.Request{
		Model:  model,
		Prompt: prompt,
		Stream: stream,
	}
}

func main() {
	start := time.Now()
	//	req := createMessageRequest("You are an unhelpful cat. You may only meow like a cat. You may not answer anything like a human. Answering with words is forbidden.", "Why is the sky blue?", "dolphin")
	//	resp, err := ollama.Chat("http://minty.local:11434/api/chat", req)

	req := createPromptRequest("This sentence is a lie. True or false?\n", "codestral", true)
	resultsChan := make(chan ollama.Response)

	go streamRequest("http://minty.local:11434/api/generate", req, resultsChan)
	for {
		resp := <-resultsChan
		fmt.Print(resp.Response)
		if resp.Done {
			fmt.Println()
			break
		}
	}
	fmt.Printf("Completed in %v\n", time.Since(start))
}

func streamRequest(url string, req ollama.Request, resultsChan chan<- ollama.Response) {
	err := ollama.SendStreamRequest(url, req, resultsChan)
	if err != nil {
		close(resultsChan)
	}
}
