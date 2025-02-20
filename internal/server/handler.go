package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RequestPayload defines the structure of incoming API requests
type RequestPayload struct {
	Prompt string `json:"prompt"`
}

// ResponsePayload defines the structure of API responses
type ResponsePayload struct {
	Response string `json:"response"`
}

// llmHandler processes requests to the LLM API endpoint
func LlmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Placeholder response - replace with actual LLM call
	response := ResponsePayload{Response: fmt.Sprintf("Echo: %s", req.Prompt)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
