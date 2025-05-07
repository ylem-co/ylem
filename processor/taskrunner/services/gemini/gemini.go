package gemini

import (
    "context"
    "fmt"
    "strings"
    "encoding/json"
    "ylem_taskrunner/config"

    "google.golang.org/genai"
)

func Process(JSON string, UserPrompt string) (string, error) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey:  config.Cfg().Gemini.Key,
        Backend: genai.BackendGeminiAPI,
    })
    if err != nil {
        return "", err
    }

    text, err := BuildPrompt(JSON, UserPrompt)
    if err != nil {
        return "", err
    }

    result, err := client.Models.GenerateContent(
        ctx,
        config.Cfg().Gemini.Model,
        genai.Text(text),
        nil,
    )
    if err != nil {
        return "", err
    }
    
    return result.Text(), nil
}

func BuildPrompt(JSON string, UserPrompt string) (string, error) {
    j, err := json.Marshal(JSON)
    if err != nil {
        return "", err
    }

    tJ := strings.Trim(string(j), "\"")

    message := fmt.Sprintf("In that JSON %s %s", tJ, UserPrompt)

    return message, nil
}
