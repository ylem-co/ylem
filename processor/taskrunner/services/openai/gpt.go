package openai

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"ylem_taskrunner/config"

	log "github.com/sirupsen/logrus"
)

type Completion struct {
	JSON       string
	UserPrompt string
}

type CompletionRequest struct {
	Model            string          `json:"model"`
	Messages         []ChoiceMessage `json:"messages"`
	Temperature      int             `json:"temperature"`
	MaxTokens        int             `json:"max_tokens"`
	TopP             int             `json:"top_p"`
	Number           int             `json:"n"`
	FrequencyPenalty int             `json:"frequency_penalty"`
	PresencePenalty  int             `json:"presence_penalty"`
}

type ChoiceMessage struct {
	Content       string `json:"content"`
	Role          string `json:"role"`
}

type ChoiceResponse struct {
	Message       ChoiceMessage `json:"message"`
	FinishReasons string                `json:"finish_reason"`
}

type CompletionResponse struct {
	Choices []ChoiceResponse `json:"choices"`
}

func (r *Completion) BuildPrompt() (*ChoiceMessage, error) {
	j, err := json.Marshal(r.JSON)
	if err != nil {
		return nil, err
	}

	tJ := strings.Trim(string(j), "\"")

	message := &ChoiceMessage{
		Role: "user",
		Content: fmt.Sprintf("In that JSON %s %s", tJ, r.UserPrompt),
	}

	return message, nil
}

func (o *OpenAi) CompleteText(c Completion) (string, error) {
	var completionResponse CompletionResponse
	message, err := c.BuildPrompt()
	if err != nil {
		return "", err
	}

	var messages []ChoiceMessage
	messages = append(messages, *message)

	request := &CompletionRequest{
		Model:            config.Cfg().Openai.Model,
		Messages:         messages,
		Temperature:      0,
		MaxTokens:        1024,
		TopP:             1,
		Number:           1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	log.Tracef("openai: gpt: complete")
	response, err := o.client.
		R().
		SetBody(request).
		SetResult(&completionResponse).
		Post("/chat/completions")

	if err != nil {
		return "", err
	}

	if response.StatusCode() != http.StatusOK {
		log.Debug(string(response.Body()))

		return "", fmt.Errorf("openai: gpt: evaluation, expected http 200, got %s", response.Status())
	}

	if len(completionResponse.Choices) == 0 {
		return "", nil
	}
	return completionResponse.Choices[0].Message.Content, nil
}
