package ollama

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages,omitempty"`
	Prompt   string    `json:"prompt,omitempty"`
	Stream   bool      `json:"stream"`
	System   string    `json:"system,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Message            Message   `json:"message"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

func SendRequest(url string, ollamaReq Request) (*Response, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	ollamaResp := Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	return &ollamaResp, err
}

func SendStreamRequest(url string, ollamaReq Request, respChan chan<- Response) error {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return err
	}

	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return err
	}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	decoder := json.NewDecoder(httpResp.Body)

	for {
		var resp Response
		if err := decoder.Decode(&resp); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		respChan <- resp

		// If it's the final response, note it and break.
		if resp.Done {
			break
		}
	}

	return nil
}
