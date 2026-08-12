package agent

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/jumpserver/wisp/pkg/agent/provider"
)

type thoughtSummaryProvider struct{}

func (thoughtSummaryProvider) Info() provider.ProviderInfo {
	return provider.ProviderInfo{}
}

func (thoughtSummaryProvider) Complete(context.Context, provider.CompletionRequest) (provider.CompletionResult, error) {
	return provider.CompletionResult{
		Content: `{
			"kind":"answer","message":"Use a read-only query","thoughtSummary":"  Chose a minimal read-only query.  ",
			"toolName":"","toolArguments":{"query":"","schema":"","tables":[],"sql":""},
			"sql":"","proposalExplanation":"",
			"analysis":{"valid":true,"statementType":"","riskLevel":0,"riskReason":"","tables":[],"columns":[],"errors":[]}
		}`,
		ReasoningContent: "private provider reasoning",
	}, nil
}

func (thoughtSummaryProvider) CompactState(provider.ContextTier) {}

type unusedSurfaceTools struct{}

func (unusedSurfaceTools) Call(context.Context, SurfaceToolCall) (json.RawMessage, error) {
	return nil, nil
}

func TestSurfaceSessionEmitsOnlyVisibleThoughtSummary(t *testing.T) {
	var messages []ChatMessage
	session := &SurfaceSession{
		sessionID: "session-1",
		provider:  thoughtSummaryProvider{},
		surface:   NewSQLSurface(),
		tools:     unusedSurfaceTools{},
		emit:      func(message ChatMessage) { messages = append(messages, message) },
	}
	session.run(context.Background(), SurfaceRequest{
		ID: "request-1", Operation: "generate", Question: "Give me a query", Context: json.RawMessage(`{}`),
	})

	var summaries []string
	for _, message := range messages {
		for _, part := range message.Parts {
			if part.Type != "data-thought-summary" {
				continue
			}
			data, ok := part.Data.(map[string]any)
			if !ok {
				t.Fatalf("thought summary data = %#v", part.Data)
			}
			summaries = append(summaries, data["text"].(string))
		}
	}
	if len(summaries) != 1 || summaries[0] != "Chose a minimal read-only query." {
		t.Fatalf("thought summaries = %#v", summaries)
	}
	for _, message := range messages {
		if strings.Contains(mustJSON(message), "private provider reasoning") {
			t.Fatal("provider reasoning content was exposed")
		}
	}
}

func TestVisibleThoughtSummaryIsBounded(t *testing.T) {
	summary := visibleThoughtSummary(strings.Repeat("界", maxSurfaceThoughtRunes+1))
	if len([]rune(summary)) != maxSurfaceThoughtRunes+1 || !strings.HasSuffix(summary, "…") {
		t.Fatalf("bounded thought summary has %d runes: %q", len([]rune(summary)), summary)
	}
}
