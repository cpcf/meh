package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const apiURL = "http://localhost:8080/api/llm"

type RequestPayload struct {
	Prompt string `json:"prompt"`
}

type ResponsePayload struct {
	Response string `json:"response"`
}

func SendRequest(prompt string) {
	data, err := json.Marshal(RequestPayload{Prompt: prompt})
	if err != nil {
		log.Fatal("Error encoding request payload:", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	var response ResponsePayload
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Error decoding response:", err)
	}

	fmt.Println(response.Response)
}

func InteractiveMode() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Interactive mode. Type your query and press Enter. Type 'exit' to quit.")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		query := scanner.Text()
		if query == "exit" {
			break
		}
		SendRequest(query)
	}
}
