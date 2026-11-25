package anthropic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"ylem_taskrunner/config"

	log "github.com/sirupsen/logrus"
)

type Completion struct {
	JSON       string
	UserPrompt string
}

type messageRequest struct {
	Model     string        `json:"model"`
	MaxTokens int           `json:"max_tokens"`
	Messages  []userMessage `json:"messages"`
}

type userMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messageResponse struct {
	Content    []contentBlock `json:"content"`
	StopReason string         `json:"stop_reason"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (c *Completion) BuildPrompt() (string, error) {
	j, err := json.Marshal(c.JSON)
	if err != nil {
		return "", err
	}

	tJ := strings.Trim(string(j), "\"")

	return fmt.Sprintf("In that JSON %s %s", tJ, c.UserPrompt), nil
}

func (a *Anthropic) CompleteText(c Completion) (string, error) {
	prompt, err := c.BuildPrompt()
	if err != nil {
		return "", err
	}

	request := &messageRequest{
		Model:     config.Cfg().Anthropic.Model,
		MaxTokens: 1024,
		Messages: []userMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	log.Tracef("anthropic: claude: complete")
	response, err := a.client.
		R().
		SetBody(request).
		Post("/messages")
	if err != nil {
		return "", err
	}

	if response.StatusCode() != http.StatusOK {
		log.Debug(string(response.Body()))

		return "", fmt.Errorf("anthropic: claude: expected http 200, got %s: %s", response.Status(), string(response.Body()))
	}

	var message messageResponse
	if err := json.Unmarshal(response.Body(), &message); err != nil {
		return "", fmt.Errorf("anthropic: claude: decode response: %w", err)
	}

	text := message.text()
	if text == "" {
		return "", fmt.Errorf("anthropic: claude: response contained no text (stop_reason=%s)", message.StopReason)
	}

	return text, nil
}

func (m messageResponse) text() string {
	var b strings.Builder
	for _, block := range m.Content {
		if block.Type == "text" {
			b.WriteString(block.Text)
		}
	}

	return b.String()
}
