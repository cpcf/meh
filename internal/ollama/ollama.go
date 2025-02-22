package ollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (m Message) String() string {
	return fmt.Sprintf("Message{Role: %s, Content: %s}", m.Role, m.Content)
}

// Request represents a call to the LLM API. It supports both chat and completion requests.
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages,omitempty"` // for chat requests
	Prompt   string    `json:"prompt,omitempty"`   // for generate (completion) requests
	System   string    `json:"system,omitempty"`   // optional: for completions with a system prompt
	Suffix   string    `json:"suffix,omitempty"`   // optional: for completions with a suffix
	Stream   bool      `json:"stream"`
}

func (r Request) String() string {
	return fmt.Sprintf("Request{Model: %s, Messages: %v, Prompt: %s, Suffix: %s, Stream: %t}", r.Model, r.Messages, r.Prompt, r.Suffix, r.Stream)
}

type Response struct {
	Model              string      `json:"model"`
	CreatedAt          time.Time   `json:"created_at"`
	Message            *Message    `json:"message,omitempty"`
	Response           string      `json:"response,omitempty"`
	Done               bool        `json:"done"`
	DoneReason         string      `json:"done_reason,omitempty"`
	TotalDuration      int64       `json:"total_duration,omitempty"`
	LoadDuration       int         `json:"load_duration,omitempty"`
	PromptEvalCount    int         `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int         `json:"prompt_eval_duration,omitempty"`
	EvalCount          int         `json:"eval_count,omitempty"`
	EvalDuration       int64       `json:"eval_duration,omitempty"`
	Context            interface{} `json:"context,omitempty"`
}

func (r Response) String() string {
	return fmt.Sprintf("Response{Model: %s, CreatedAt: %s, Message: %v, Response: %s, Done: %t, DoneReason: %s, TotalDuration: %d, LoadDuration: %d, PromptEvalCount: %d, PromptEvalDuration: %d, EvalCount: %d, EvalDuration: %d, Context: %v}", r.Model, r.CreatedAt, r.Message, r.Response, r.Done, r.DoneReason, r.TotalDuration, r.LoadDuration, r.PromptEvalCount, r.PromptEvalDuration, r.EvalCount, r.EvalDuration, r.Context)
}

// SendRequest sends a non-streaming HTTP POST request to the given endpoint.
func sendRequest(url string, req Request) (*Response, error) {
	js, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpResp, err := http.Post(url, "application/json", bytes.NewReader(js))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", httpResp.Status)
	}

	var apiResp Response
	if err := json.NewDecoder(httpResp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &apiResp, nil
}

// sendStreamRequest sends a streaming HTTP POST request to the given endpoint.
// It decodes a series of JSON responses and writes each to respChan.
func sendStreamRequest(url string, req Request, respChan chan<- Response) error {
	js, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpResp, err := http.Post(url, "application/json", bytes.NewReader(js))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	decoder := json.NewDecoder(httpResp.Body)
	for {
		var resp Response
		if err := decoder.Decode(&resp); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}
		respChan <- resp
		if resp.Done {
			break
		}
	}

	return nil
}

// runningModelsResponse is used to decode the JSON response from GET /tags.
type runningModelsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

type OllamaAPI struct {
	baseURL      string
	model        string
	systemPrompt string
	history      []Message
	closed       bool
}

func NewAPI(baseURL, system string) *OllamaAPI {
	// if we have a system prompt, add it to the start of the history
	history := make([]Message, 0)
	if system != "" {
		history = append(history, Message{Role: "system", Content: system})
	}

	return &OllamaAPI{
		baseURL:      baseURL,
		systemPrompt: system,
		history:      history,
	}

}

// Chat sends a chat message using the /chat endpoint.
// It records the conversation history, sends all messages on each call,
// and appends the assistant’s reply to the history.
func (o *OllamaAPI) Chat(message string, results chan string, stream bool) {
	if o.closed {
		results <- "Error: API is closed"
		return
	}
	// Append the user's message.
	o.history = append(o.history, Message{Role: "user", Content: message})
	req := Request{
		Model:    o.model,
		Messages: o.history,
		Stream:   stream,
	}

	endpoint := o.baseURL + "/chat"

	if stream {
		respChan := make(chan Response)
		go func() {
			if err := sendStreamRequest(endpoint, req, respChan); err != nil {
				results <- fmt.Sprintf("Error: %v", err)
			}
			close(respChan)
		}()

		var fullResponse string
		for resp := range respChan {
			results <- resp.Message.Content
			fullResponse += resp.Message.Content
			if resp.Done {
				break
			}
		}
		close(results)
		// Append assistant's reply.
		o.history = append(o.history, Message{Role: "assistant", Content: fullResponse})
	} else {
		resp, err := sendRequest(endpoint, req)
		if err != nil {
			results <- fmt.Sprintf("Error: %v", err)
			return
		}
		results <- resp.Message.Content
		o.history = append(o.history, Message{Role: "assistant", Content: resp.Response})
	}
}

// Prompt sends a prompt using the /generate endpoint (completion).
func (o *OllamaAPI) Prompt(message string, results chan string, stream bool) {
	if o.closed {
		results <- "Error: API is closed"
		return
	}

	req := Request{
		Model:  o.model,
		Prompt: message,
		Stream: stream,
	}

	if o.systemPrompt != "" {
		req.System = o.systemPrompt
	}

	endpoint := o.baseURL + "/generate"

	if stream {
		respChan := make(chan Response)
		go func() {
			if err := sendStreamRequest(endpoint, req, respChan); err != nil {
				results <- fmt.Sprintf("Error: %v", err)
			}
			close(respChan)
		}()
		for resp := range respChan {
			results <- resp.Response
			if resp.Done {
				break
			}
		}
		close(results)

	} else {
		resp, err := sendRequest(endpoint, req)
		if err != nil {
			results <- fmt.Sprintf("Error: %v", err)
			return
		}
		results <- resp.Response
	}
}

// Models retrieves the list of local models using the /tags endpoint.
func (o *OllamaAPI) Models() []string {
	if o.closed {
		return []string{}
	}

	endpoint := o.baseURL + "/tags"

	httpResp, err := http.Get(endpoint)
	if err != nil {
		fmt.Printf("Error fetching models: %v\n", err)
		return []string{}
	}
	defer httpResp.Body.Close()

	var resp runningModelsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		fmt.Printf("Error decoding models response: %v\n", err)
		return []string{}
	}

	models := make([]string, 0, len(resp.Models))
	for _, m := range resp.Models {
		models = append(models, m.Name)
	}
	return models
}

func (o *OllamaAPI) SelectModel(model string) {
	o.model = model
}
