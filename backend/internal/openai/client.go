package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultModel   = "gpt-4o"
	defaultTimeout = 120 * time.Second
)

var (
	ErrEmptyAPIKey     = errors.New("OPENAI_API_KEY environment variable is not set")
	ErrEmptyInput      = errors.New("input cannot be empty or whitespace only")
	ErrRequestFailed   = errors.New("openai request failed")
	ErrInvalidResponse = errors.New("invalid response from openai")
)

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents the request body for chat completions.
type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatCompletionResponse represents the response from chat completions.
type ChatCompletionResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice represents a single completion choice.
type Choice struct {
	Message Message `json:"message"`
}

// APIError represents an error from the OpenAI API.
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Client is an OpenAI API client.
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
	model      string
}

// NewClient creates a new OpenAI client.
// It loads the API key from the OPENAI_API_KEY environment variable.
func NewClient() (*Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, ErrEmptyAPIKey
	}

	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: defaultBaseURL,
		model:   defaultModel,
	}, nil
}

// NewClientWithConfig creates a new OpenAI client with custom configuration.
func NewClientWithConfig(apiKey, baseURL, model string, timeout time.Duration) (*Client, error) {
	if apiKey == "" {
		return nil, ErrEmptyAPIKey
	}

	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if model == "" {
		model = defaultModel
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
		model:   model,
	}, nil
}

// ValidateInput checks if the input is valid (non-empty and not whitespace only).
func ValidateInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return ErrEmptyInput
	}
	return nil
}

// ChatCompletion sends a chat completion request to the OpenAI API.
// The context can be used to set a timeout or cancel the request.
func (c *Client) ChatCompletion(ctx context.Context, messages []Message) (string, error) {
	if len(messages) == 0 {
		return "", ErrEmptyInput
	}

	reqBody := ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("request timed out: %w", err)
		}
		return "", fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ChatCompletionResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != nil {
			return "", fmt.Errorf("%w: %s", ErrRequestFailed, errResp.Error.Message)
		}
		return "", fmt.Errorf("%w: status %d", ErrRequestFailed, resp.StatusCode)
	}

	var completionResp ChatCompletionResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("%w: no choices returned", ErrInvalidResponse)
	}

	return completionResp.Choices[0].Message.Content, nil
}
