package perplexity

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiURL    = "https://api.perplexity.ai/chat/completions"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
)

const (
	CLAUDE = "mixtral-8x7b-instruct"  // Non-online model that might be available
)

var (
	// Make sure we don't get cloudflare'd.
	tlsConfig = &tls.Config{
		MinVersion:       tls.VersionTLS13,
		CipherSuites:     []uint16{tls.TLS_AES_128_GCM_SHA256},
		CurvePreferences: []tls.CurveID{tls.X25519},
	}

	ALL_MODELS = []string{
		"mixtral-8x7b-instruct",
		"pplx-7b-chat",
		"pplx-70b-chat",
	}
)

type ChatClient struct {
	client           *http.Client
	History          []Message
	debug            bool
	conversationMode bool
	apiKey           string
}

// Add a struct for the API response
type APIResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int     `json:"index"`
		Message Message `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewChatClient(debug bool, conversationMode bool) *ChatClient {
	return &ChatClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
			Timeout: time.Second * 60,
		},
		History:          []Message{},
		debug:            debug,
		conversationMode: conversationMode,
		apiKey:           "", // Will need to set this via environment variable or config
	}
}

// SetAPIKey allows setting the API key externally
func (c *ChatClient) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *ChatClient) Backtrack() error {
	if len(c.History) < 2 {
		return errors.New("no history to backtrack")
	}
	c.History = c.History[:len(c.History)-2]
	return nil
}

func (c *ChatClient) Connect() error {
	// Check if API key is set
	if c.apiKey == "" {
		return errors.New("API key not set. Please set your Perplexity API key using PPLX_API_KEY environment variable or via SetAPIKey method")
	}
	return nil
}

func (c *ChatClient) Close() error {
	// No explicit close needed for HTTP client
	return nil
}

func (c *ChatClient) ReceiveMessage(timeout time.Duration) (string, error) {
	// This method doesn't make sense in the new HTTP-based API,
	// but we keep it for compatibility with the existing interface
	// In a real implementation, responses are immediate with HTTP
	return "", errors.New("ReceiveMessage not supported in HTTP API mode. Use SendMessage which returns the response directly")
}

func (c *ChatClient) SendMessage(message string, model string) error {
	if c.conversationMode {
		fmt.Println(fmt.Sprintf("\r\nUser:\r\n %s\r\n", message))
	}

	// Add user message to history
	userMessage := Message{
		Role:    "user",
		Content: message,
	}
	c.History = append(c.History, userMessage)

	// Prepare request payload
	req := Request{
		Model:    model,
		Messages: c.History,
		Stream:   false, // For simplicity, non-streaming for now
	}

	// Convert to the format expected by the API
	apiReq := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
	}

	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return err
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("User-Agent", userAgent)

	// Make the API call
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var apiResponse APIResponse
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return fmt.Errorf("error unmarshaling API response: %v", err)
	}

	if len(apiResponse.Choices) == 0 {
		return errors.New("no choices returned in API response")
	}

	// Get the assistant's response
	assistantResponse := apiResponse.Choices[0].Message

	// Add assistant response to history
	c.History = append(c.History, assistantResponse)

	if c.conversationMode {
		fmt.Println(fmt.Sprintf("\r\nAssistant:\r\n %s\r\n", assistantResponse.Content))
	}

	return nil
}

// New method to get response directly instead of using ReceiveMessage
func (c *ChatClient) GetLastResponse() string {
	if len(c.History) > 0 {
		// Return the last message which should be the assistant's response
		return c.History[len(c.History)-1].Content
	}
	return ""
}

// Updated Request struct to match OpenAI/Perplexity API format
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
	MaxTokens *int     `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	TopP *float64      `json:"top_p,omitempty"`
}
