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

type repeatingSQLToolProvider struct {
	calls    int
	requests []provider.CompletionRequest
}

func (p *repeatingSQLToolProvider) Info() provider.ProviderInfo {
	return provider.ProviderInfo{}
}

func (p *repeatingSQLToolProvider) Complete(
	_ context.Context,
	request provider.CompletionRequest,
) (provider.CompletionResult, error) {
	p.calls++
	p.requests = append(p.requests, request)
	return provider.CompletionResult{Content: `{
		"kind":"tool","message":"","thoughtSummary":"Inspecting users",
		"toolName":"inspect_schema",
		"toolArguments":{"query":"users","schema":"public","tables":[],"sql":""},
		"sql":"","proposalExplanation":"",
		"analysis":{"valid":false,"statementType":"","riskLevel":0,
		"riskReason":"","tables":[],"columns":[],"errors":[]}
	}`}, nil
}

func (p *repeatingSQLToolProvider) CompactState(provider.ContextTier) {}

type countingSurfaceTools struct {
	calls int
}

func (t *countingSurfaceTools) Call(context.Context, SurfaceToolCall) (json.RawMessage, error) {
	t.calls++
	return json.RawMessage(`{"tables":[]}`), nil
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

func TestSurfaceSessionStopsRepeatedSQLToolLoop(t *testing.T) {
	model := &repeatingSQLToolProvider{}
	tools := &countingSurfaceTools{}
	surface := NewSQLSurface()
	surface.SetLanguage("zh-CN")
	var messages []ChatMessage
	session := &SurfaceSession{
		sessionID: "session-loop",
		provider:  model,
		surface:   surface,
		tools:     tools,
		emit:      func(message ChatMessage) { messages = append(messages, message) },
	}
	session.run(context.Background(), SurfaceRequest{
		ID: "request-loop", Operation: "generate", Question: "查询用户",
		Context: json.RawMessage(`{}`),
	})

	if model.calls != 3 {
		t.Fatalf("model calls = %d, want 3", model.calls)
	}
	if tools.calls != 1 {
		t.Fatalf("tool calls = %d, want one real call", tools.calls)
	}
	if len(model.requests) != 3 {
		t.Fatalf("captured model requests = %d", len(model.requests))
	}
	properties := model.requests[2].Tool.Parameters["properties"].(map[string]any)
	kinds := properties["kind"].(map[string]any)["enum"].([]string)
	if len(kinds) != 2 || kinds[0] != "answer" || kinds[1] != "proposal" {
		t.Fatalf("final model request kinds = %#v", kinds)
	}

	var finalText string
	var timing map[string]any
	for _, message := range messages {
		for _, part := range message.Parts {
			if part.Type == "data-error" {
				t.Fatalf("loop protection emitted data-error: %#v", part.Data)
			}
			if part.Type == "text" && message.Metadata["stage"] == "final" {
				finalText = part.Text
			}
			if part.Type == "data-agent-timing" {
				timing = part.Data.(map[string]any)
			}
		}
	}
	if finalText != "我无法从当前数据库结构中确定所需的对象。这段 SQL 应使用哪些表和字段？" {
		t.Fatalf("final fallback text = %q", finalText)
	}
	if timing == nil || timing["toolCalls"] != 1 ||
		timing["duplicateToolBlocked"] != 1 ||
		timing["schemaBudgetExhausted"] != 0 ||
		timing["forcedClarifications"] != 1 {
		t.Fatalf("loop protection timing = %#v", timing)
	}
}
