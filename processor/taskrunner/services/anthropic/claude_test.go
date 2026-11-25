package anthropic

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestBuildPrompt(t *testing.T) {
	c := Completion{
		JSON:       `{"amount":1}`,
		UserPrompt: "Multiply all amounts by 3 and return new JSON only.",
	}

	got, err := c.BuildPrompt()
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
	}

	want := `In that JSON {\"amount\":1} Multiply all amounts by 3 and return new JSON only.`
	if got != want {
		t.Fatalf("BuildPrompt() = %q, want %q", got, want)
	}
}

func TestCompleteText(t *testing.T) {
	var gotBody messageRequest
	var gotAuth string
	var gotVersion string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			t.Errorf("path = %s, want /messages", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		gotAuth = r.Header.Get("Authorization")
		gotVersion = r.Header.Get("anthropic-version")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		if err := json.Unmarshal(body, &gotBody); err != nil {
			t.Errorf("decode body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{
			"content": [
				{"type": "thinking", "thinking": "ignore me"},
				{"type": "text", "text": "{\"amount\":3}"}
			],
			"stop_reason": "end_turn"
		}`)
	}))
	defer srv.Close()

	client := &Anthropic{
		client: resty.New().
			SetBaseURL(srv.URL).
			SetAuthToken("test-key").
			SetHeader("anthropic-version", apiVersion),
	}

	got, err := client.CompleteText(Completion{
		JSON:       `{"amount":1}`,
		UserPrompt: "Multiply all amounts by 3",
	})
	if err != nil {
		t.Fatalf("CompleteText() error = %v", err)
	}
	if got != `{"amount":3}` {
		t.Fatalf("CompleteText() = %q, want JSON text only", got)
	}
	if gotAuth != "Bearer test-key" {
		t.Fatalf("Authorization = %q, want bearer token", gotAuth)
	}
	if gotVersion != apiVersion {
		t.Fatalf("anthropic-version = %q, want %q", gotVersion, apiVersion)
	}
	if gotBody.MaxTokens != 1024 {
		t.Fatalf("max_tokens = %d, want 1024", gotBody.MaxTokens)
	}
	if len(gotBody.Messages) != 1 || gotBody.Messages[0].Role != "user" {
		t.Fatalf("messages = %#v, want one user message", gotBody.Messages)
	}
	if !strings.Contains(gotBody.Messages[0].Content, "Multiply all amounts by 3") {
		t.Fatalf("prompt = %q, want the user prompt", gotBody.Messages[0].Content)
	}
}

func TestCompleteTextHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"type":"error","error":{"message":"invalid api key"}}`, http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := &Anthropic{
		client: resty.New().SetBaseURL(srv.URL),
	}

	_, err := client.CompleteText(Completion{JSON: `{}`, UserPrompt: "summarize"})
	if err == nil {
		t.Fatal("CompleteText() error = nil, want status error")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Fatalf("CompleteText() error = %v, want http 401", err)
	}
}

func TestCompleteTextNoText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"content":[{"type":"thinking","thinking":"only"}],"stop_reason":"refusal"}`)
	}))
	defer srv.Close()

	client := &Anthropic{
		client: resty.New().SetBaseURL(srv.URL),
	}

	_, err := client.CompleteText(Completion{JSON: `{}`, UserPrompt: "summarize"})
	if err == nil {
		t.Fatal("CompleteText() error = nil, want empty-text error")
	}
	if !strings.Contains(err.Error(), "no text") || !strings.Contains(err.Error(), "refusal") {
		t.Fatalf("CompleteText() error = %v, want no-text refusal", err)
	}
}
