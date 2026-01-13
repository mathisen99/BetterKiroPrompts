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
	defaultModel   = "gpt-5.2"
	defaultTimeout = 120 * time.Second
)

// ReasoningEffort controls how many reasoning tokens the model generates.
type ReasoningEffort string

const (
	ReasoningNone   ReasoningEffort = "none"
	ReasoningLow    ReasoningEffort = "low"
	ReasoningMedium ReasoningEffort = "medium"
	ReasoningHigh   ReasoningEffort = "high"
	ReasoningXHigh  ReasoningEffort = "xhigh"
)

// Verbosity controls output token generation.
type Verbosity string

const (
	VerbosityLow    Verbosity = "low"
	VerbosityMedium Verbosity = "medium"
	VerbosityHigh   Verbosity = "high"
)

var (
	ErrEmptyAPIKey     = errors.New("OPENAI_API_KEY environment variable is not set")
	ErrEmptyInput      = errors.New("input cannot be empty or whitespace only")
	ErrRequestFailed   = errors.New("openai request failed")
	ErrInvalidResponse = errors.New("invalid response from openai")
)

// Message represents a chat message (used for building input).
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Reasoning configures reasoning behavior for GPT-5.2.
type Reasoning struct {
	Effort ReasoningEffort `json:"effort,omitempty"`
}

// TextConfig configures text output behavior.
type TextConfig struct {
	Verbosity Verbosity `json:"verbosity,omitempty"`
}

// ResponsesRequest represents the request body for the Responses API.
type ResponsesRequest struct {
	Model              string      `json:"model"`
	Input              any         `json:"input"`
	Reasoning          *Reasoning  `json:"reasoning,omitempty"`
	Text               *TextConfig `json:"text,omitempty"`
	PreviousResponseID string      `json:"previous_response_id,omitempty"`
}

// ResponsesResponse represents the response from the Responses API.
type ResponsesResponse struct {
	ID         string       `json:"id"`
	Output     []OutputItem `json:"output"`
	OutputText string       `json:"output_text"` // Convenience field aggregating all text
	Error      *APIError    `json:"error,omitempty"`
}

// OutputItem represents an item in the response output array.
type OutputItem struct {
	Type    string         `json:"type"`
	ID      string         `json:"id,omitempty"`
	Role    string         `json:"role,omitempty"`
	Content []ContentBlock `json:"content,omitempty"`
}

// ContentBlock represents a content block in the output.
type ContentBlock struct {
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"`
	Annotations []any  `json:"annotations,omitempty"`
}

// APIError represents an error from the OpenAI API.
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Client is an OpenAI API client configured for GPT-5.2.
type Client struct {
	apiKey          string
	httpClient      *http.Client
	baseURL         string
	model           string
	reasoningEffort ReasoningEffort
	verbosity       Verbosity
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
		baseURL:         defaultBaseURL,
		model:           defaultModel,
		reasoningEffort: ReasoningMedium,
		verbosity:       VerbosityMedium,
	}, nil
}

// ClientConfig holds configuration options for the client.
type ClientConfig struct {
	APIKey          string
	BaseURL         string
	Model           string
	Timeout         time.Duration
	ReasoningEffort ReasoningEffort
	Verbosity       Verbosity
}

// NewClientWithConfig creates a new OpenAI client with custom configuration.
func NewClientWithConfig(cfg ClientConfig) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.Model == "" {
		cfg.Model = defaultModel
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}
	if cfg.ReasoningEffort == "" {
		cfg.ReasoningEffort = ReasoningMedium
	}
	if cfg.Verbosity == "" {
		cfg.Verbosity = VerbosityMedium
	}

	return &Client{
		apiKey: cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		baseURL:         cfg.BaseURL,
		model:           cfg.Model,
		reasoningEffort: cfg.ReasoningEffort,
		verbosity:       cfg.Verbosity,
	}, nil
}

// SetReasoningEffort updates the reasoning effort level.
func (c *Client) SetReasoningEffort(effort ReasoningEffort) {
	c.reasoningEffort = effort
}

// SetVerbosity updates the verbosity level.
func (c *Client) SetVerbosity(v Verbosity) {
	c.verbosity = v
}

// ValidateInput checks if the input is valid (non-empty and not whitespace only).
func ValidateInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return ErrEmptyInput
	}
	return nil
}

// ChatCompletion sends a request to the GPT-5.2 Responses API.
// The context can be used to set a timeout or cancel the request.
func (c *Client) ChatCompletion(ctx context.Context, messages []Message) (string, error) {
	if len(messages) == 0 {
		return "", ErrEmptyInput
	}

	// Convert messages to Responses API input format
	input := convertMessagesToInput(messages)

	reqBody := ResponsesRequest{
		Model: c.model,
		Input: input,
		Reasoning: &Reasoning{
			Effort: c.reasoningEffort,
		},
		Text: &TextConfig{
			Verbosity: c.verbosity,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/responses", bytes.NewReader(jsonBody))
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
		var errResp ResponsesResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != nil {
			return "", fmt.Errorf("%w: %s", ErrRequestFailed, errResp.Error.Message)
		}
		return "", fmt.Errorf("%w: status %d: %s", ErrRequestFailed, resp.StatusCode, string(body))
	}

	var responsesResp ResponsesResponse
	if err := json.Unmarshal(body, &responsesResp); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}

	// Use the convenience output_text field if available
	if responsesResp.OutputText != "" {
		return responsesResp.OutputText, nil
	}

	// Fall back to extracting from output array
	text := extractTextFromResponse(responsesResp)
	if text == "" {
		return "", fmt.Errorf("%w: no text content in response", ErrInvalidResponse)
	}

	return text, nil
}

// convertMessagesToInput converts Message slice to Responses API input format.
// Maps "system" role to "developer" for GPT-5.2 compatibility.
func convertMessagesToInput(messages []Message) []map[string]any {
	input := make([]map[string]any, len(messages))
	for i, msg := range messages {
		role := msg.Role
		// GPT-5.2 uses "developer" instead of "system" for instruction messages
		if role == "system" {
			role = "developer"
		}
		input[i] = map[string]any{
			"role":    role,
			"content": msg.Content,
		}
	}
	return input
}

// extractTextFromResponse extracts the text content from a Responses API response.
func extractTextFromResponse(resp ResponsesResponse) string {
	var texts []string
	for _, item := range resp.Output {
		// Handle message type output items
		if item.Type == "message" {
			for _, block := range item.Content {
				if block.Type == "output_text" && block.Text != "" {
					texts = append(texts, block.Text)
				}
			}
		}
	}
	return strings.Join(texts, "")
}
