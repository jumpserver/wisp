package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jumpserver/wisp/pkg/agent/provider"
	"github.com/jumpserver/wisp/pkg/logger"
)

const (
	maxSurfaceQuestionBytes   = 32 * 1024
	maxSurfaceContextBytes    = 256 * 1024
	maxSurfaceToolCalls       = 12
	maxSurfaceToolArguments   = 32 * 1024
	maxSurfaceToolResultBytes = 256 * 1024
	maxSurfaceRounds          = 16
	maxSurfaceHistoryBytes    = 1024 * 1024
	maxSurfaceThoughtRunes    = 400
	defaultModelRequestLimit  = 30
)

type SurfaceRequest struct {
	ID        string          `json:"id"`
	Operation string          `json:"operation"`
	Question  string          `json:"question"`
	Context   json.RawMessage `json:"context"`
}

type SurfaceToolCall struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type SurfaceToolResult struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     string          `json:"error,omitempty"`
}

type SurfaceState struct {
	History           string              `json:"history,omitempty"`
	ToolResults       []SurfaceToolResult `json:"toolResults,omitempty"`
	Correction        string              `json:"correction,omitempty"`
	Round             int                 `json:"round"`
	MaximumRound      int                 `json:"maximumRound"`
	ToolCallsDisabled bool                `json:"-"`
}

type SurfaceAction struct {
	Kind        string
	Text        string
	Thought     string
	Tool        *SurfaceToolCall
	Value       any
	HistoryText string
}

type SurfaceReview struct {
	Tool              *SurfaceToolCall
	Correction        string
	FinalizeAfterTool bool
}

// Surface supplies a model contract and policy for one kind of interactive
// workspace. The runtime owns lifecycle, budgets, tool correlation and audit;
// the surface owns domain prompts and final UI parts.
type Surface interface {
	Name() string
	CompletionRequest(SurfaceRequest, SurfaceState, provider.ContextTier) provider.CompletionRequest
	DecodeAction(string) (SurfaceAction, error)
	ValidateTool(SurfaceToolCall) error
	Review(SurfaceRequest, SurfaceState, SurfaceAction) (SurfaceReview, error)
	FinalParts(SurfaceRequest, SurfaceState, SurfaceAction) ([]ChatPart, error)
}

type SurfaceToolCaller interface {
	Call(context.Context, SurfaceToolCall) (json.RawMessage, error)
}

type SurfaceInitializer interface {
	InitialTools(SurfaceRequest) ([]SurfaceToolCall, error)
}

type SurfaceToolCallPolicyResult struct {
	Blocked             bool
	DisableFurtherTools bool
	Correction          string
	Outcome             string
}

// SurfaceToolCallPolicy lets a surface stop domain-specific no-progress tool
// loops before they consume the shared runtime tool-call budget.
type SurfaceToolCallPolicy interface {
	EvaluateToolCall(SurfaceRequest, SurfaceState, SurfaceToolCall) SurfaceToolCallPolicyResult
}

// SurfaceToolResultPolicy lets a surface stop retrying a tool after a result
// proves that further calls cannot make progress for the current request.
type SurfaceToolResultPolicy interface {
	EvaluateToolResult(SurfaceState, SurfaceToolResult) SurfaceToolCallPolicyResult
}

// SurfaceToolCallFallback supplies final user-visible content when a model
// ignores a tool-disabled correction or returns invalid structured output.
type SurfaceToolCallFallback interface {
	ToolCallsDisabledAction(SurfaceRequest, SurfaceState) (SurfaceAction, error)
}

type surfaceRequestMetrics struct {
	started               time.Time
	rounds                int
	modelRequests         int
	toolCalls             int
	duplicateToolBlocked  int
	schemaBudgetExhausted int
	forcedClarifications  int
	modelDuration         time.Duration
	toolDuration          time.Duration
	queueDuration         time.Duration
}

func newSurfaceRequestMetrics() *surfaceRequestMetrics {
	return &surfaceRequestMetrics{started: time.Now()}
}

func (m *surfaceRequestMetrics) data() map[string]any {
	if m == nil {
		return map[string]any{}
	}
	return map[string]any{
		"durationMs":            durationMilliseconds(time.Since(m.started)),
		"rounds":                m.rounds,
		"modelRequests":         m.modelRequests,
		"modelDurationMs":       durationMilliseconds(m.modelDuration),
		"toolCalls":             m.toolCalls,
		"duplicateToolBlocked":  m.duplicateToolBlocked,
		"schemaBudgetExhausted": m.schemaBudgetExhausted,
		"forcedClarifications":  m.forcedClarifications,
		"toolDurationMs":        durationMilliseconds(m.toolDuration),
		"queueDurationMs":       durationMilliseconds(m.queueDuration),
	}
}

func durationMilliseconds(value time.Duration) float64 {
	return float64(value.Microseconds()) / 1000
}

type SurfaceSessionOptions struct {
	SessionID      string
	UserID         string
	Language       string
	Config         Config
	Surface        Surface
	Tools          SurfaceToolCaller
	AcquireRequest func(context.Context) (func(), error)
	Emit           func(ChatMessage)
}

type SurfaceSession struct {
	sessionID      string
	provider       provider.Provider
	config         Config
	surface        Surface
	tools          SurfaceToolCaller
	acquireRequest func(context.Context) (func(), error)
	emit           func(ChatMessage)
	audit          *auditWriter

	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc

	mu      sync.Mutex
	busy    bool
	closed  bool
	cancel  context.CancelFunc
	history []string
	wg      sync.WaitGroup
}

func NewSurfaceSession(options SurfaceSessionOptions) (*SurfaceSession, error) {
	if strings.TrimSpace(options.SessionID) == "" {
		return nil, fmt.Errorf("agent surface session id is required")
	}
	if options.Surface == nil || options.Tools == nil || options.Emit == nil {
		return nil, fmt.Errorf("agent surface dependencies are incomplete")
	}
	audit := newAuditWriter(
		options.UserID, options.Config.MemoryRoot, options.Config.MemorySessions,
	)
	config := options.Config
	config.Provider.Trace = audit
	modelProvider, err := provider.New(config.Provider)
	if err != nil {
		if audit != nil {
			audit.Close()
		}
		return nil, err
	}
	lifetimeCtx, lifetimeCancel := context.WithCancel(context.Background())
	session := &SurfaceSession{
		sessionID:      options.SessionID,
		provider:       modelProvider,
		config:         config,
		surface:        options.Surface,
		tools:          options.Tools,
		acquireRequest: options.AcquireRequest,
		emit:           options.Emit,
		audit:          audit,
		lifetimeCtx:    lifetimeCtx, lifetimeCancel: lifetimeCancel,
	}
	if audit != nil {
		audit.SetSessionID(options.SessionID)
	}
	if languageAware, ok := options.Surface.(interface{ SetLanguage(string) }); ok {
		languageAware.SetLanguage(options.Language)
	}
	return session, nil
}

func (s *SurfaceSession) ProviderInfo() provider.ProviderInfo {
	return s.provider.Info()
}

func (s *SurfaceSession) AnnounceCapability() {
	info := s.provider.Info()
	s.emitData("data-capability", map[string]any{
		"enabled": true, "surface": s.surface.Name(),
		"provider": info.Name, "model": info.Model,
		"modelCapabilities": info.Capabilities,
		"executionEnabled":  false,
	}, "process", "")
}

func (s *SurfaceSession) Handle(request SurfaceRequest) error {
	request.ID = strings.TrimSpace(request.ID)
	request.Operation = strings.ToLower(strings.TrimSpace(request.Operation))
	request.Question = strings.TrimSpace(request.Question)
	if request.ID == "" || request.Question == "" {
		return fmt.Errorf("agent surface request id and question are required")
	}
	if len(request.Question) > maxSurfaceQuestionBytes {
		return fmt.Errorf("agent surface question is too large")
	}
	if len(request.Context) == 0 || len(request.Context) > maxSurfaceContextBytes || !json.Valid(request.Context) {
		return fmt.Errorf("agent surface context is invalid")
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return fmt.Errorf("agent surface session is closed")
	}
	if s.busy {
		s.mu.Unlock()
		return fmt.Errorf("another agent request is active")
	}
	ctx, cancel := context.WithCancel(s.lifetimeCtx)
	s.busy = true
	s.cancel = cancel
	s.history = append(s.history, "user: "+request.Question)
	s.trimHistoryLocked()
	s.mu.Unlock()

	s.writeAudit("surface_request", request)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer cancel()
		s.run(ctx, request)
		s.mu.Lock()
		s.cancel = nil
		s.busy = false
		s.mu.Unlock()
	}()
	return nil
}

func (s *SurfaceSession) run(ctx context.Context, request SurfaceRequest) {
	metrics := newSurfaceRequestMetrics()
	lastThought := ""
	s.emitProgress("Analyzing request", "analyzing", "running", true, request.ID, metrics, nil)
	state := SurfaceState{MaximumRound: maxSurfaceRounds}
	requestBudget := s.config.MaxModelRequests
	if requestBudget <= 0 {
		requestBudget = defaultModelRequestLimit
	}
	ctx = provider.WithRequestBudget(ctx, requestBudget)
	if initializer, ok := s.surface.(SurfaceInitializer); ok {
		calls, err := initializer.InitialTools(request)
		if err != nil {
			s.finishWithError(request.ID, err, metrics)
			return
		}
		for _, call := range calls {
			if err = s.callTool(ctx, request.ID, &state, call, metrics); err != nil {
				s.finishWithError(request.ID, err, metrics)
				return
			}
		}
	}

	for round := 1; round <= maxSurfaceRounds; round++ {
		state.Round = round
		metrics.rounds = round
		s.mu.Lock()
		state.History = headTailPrompt(strings.Join(s.history, "\n"), maxSurfaceHistoryBytes)
		s.mu.Unlock()

		s.emitProgress("Waiting for AI model", "model", "running", true, request.ID, metrics, map[string]any{
			"round": round, "maximumRound": maxSurfaceRounds,
			"modelRequests": metrics.modelRequests + 1,
		})
		content, err := s.complete(ctx, request, state, metrics)
		if err != nil {
			s.finishWithError(request.ID, err, metrics)
			return
		}
		action, err := s.surface.DecodeAction(content)
		if err != nil {
			if state.ToolCallsDisabled {
				action, err = s.toolCallsDisabledAction(request, state, metrics, "invalid_output")
				if err != nil {
					s.finishWithError(request.ID, err, metrics)
					return
				}
			} else {
				state.Correction = err.Error()
				s.writeAudit("surface_output_repair", map[string]any{
					"requestId": request.ID, "round": round, "error": err.Error(),
				})
				continue
			}
		}
		state.Correction = ""
		if action.Tool != nil {
			if state.ToolCallsDisabled {
				action, err = s.toolCallsDisabledAction(request, state, metrics, "tool_requested_while_disabled")
				if err != nil {
					s.finishWithError(request.ID, err, metrics)
					return
				}
			} else if policy, ok := s.surface.(SurfaceToolCallPolicy); ok {
				decision := policy.EvaluateToolCall(request, state, *action.Tool)
				if decision.Blocked {
					state.Correction = strings.TrimSpace(decision.Correction)
					state.ToolCallsDisabled = decision.DisableFurtherTools
					s.recordToolPolicy(request.ID, decision.Outcome, *action.Tool, metrics)
					continue
				}
			}
		}
		thought := visibleThoughtSummary(action.Thought)
		if thought != "" && thought != lastThought {
			s.emitData("data-thought-summary", map[string]any{"text": thought}, "process", request.ID)
			lastThought = thought
		}
		if action.Tool != nil {
			if err = s.callTool(ctx, request.ID, &state, *action.Tool, metrics); err != nil {
				s.finishWithError(request.ID, err, metrics)
				return
			}
			continue
		}

		review, err := s.surface.Review(request, state, action)
		if err != nil {
			s.finishWithError(request.ID, err, metrics)
			return
		}
		if review.Tool != nil {
			finalizeAfterTool := review.FinalizeAfterTool
			if err = s.callTool(ctx, request.ID, &state, *review.Tool, metrics); err != nil {
				s.finishWithError(request.ID, err, metrics)
				return
			}
			if !finalizeAfterTool {
				continue
			}
			review, err = s.surface.Review(request, state, action)
			if err != nil {
				s.finishWithError(request.ID, err, metrics)
				return
			}
			if review.Tool != nil {
				continue
			}
		}
		if review.Correction != "" {
			state.Correction = review.Correction
			continue
		}

		parts, err := s.surface.FinalParts(request, state, action)
		if err != nil {
			s.finishWithError(request.ID, err, metrics)
			return
		}
		if len(parts) == 0 {
			s.finishWithError(request.ID, fmt.Errorf("agent surface returned no final content"), metrics)
			return
		}
		parts = append(parts, ChatPart{Type: "data-agent-timing", Data: metrics.data()})
		s.emitMessage(parts, "final", request.ID)
		history := strings.TrimSpace(action.HistoryText)
		if history == "" {
			history = strings.TrimSpace(action.Text)
		}
		if history != "" {
			s.mu.Lock()
			s.history = append(s.history, "assistant: "+headTailPrompt(history, 32*1024))
			s.trimHistoryLocked()
			s.mu.Unlock()
		}
		s.writeAudit("surface_final", map[string]any{
			"requestId": request.ID, "kind": action.Kind, "parts": parts,
		})
		s.logTiming(request.ID, "complete", metrics)
		s.emitProgress("", "complete", "idle", false, request.ID, metrics, nil)
		return
	}
	s.finishWithError(request.ID, fmt.Errorf("agent surface reached the maximum reasoning rounds"), metrics)
}

func (s *SurfaceSession) toolCallsDisabledAction(
	request SurfaceRequest,
	state SurfaceState,
	metrics *surfaceRequestMetrics,
	reason string,
) (SurfaceAction, error) {
	fallback, ok := s.surface.(SurfaceToolCallFallback)
	if !ok {
		return SurfaceAction{}, fmt.Errorf("agent surface requested a tool after tools were disabled")
	}
	action, err := fallback.ToolCallsDisabledAction(request, state)
	if err != nil {
		return SurfaceAction{}, err
	}
	metrics.forcedClarifications++
	logger.Infof(
		"Agent timing surface=%s request=%s stage=tool_policy outcome=forced_clarification reason=%s",
		s.surface.Name(), request.ID, reason,
	)
	s.writeAudit("surface_tool_policy", map[string]any{
		"requestId": request.ID, "outcome": "forced_clarification", "reason": reason,
	})
	return action, nil
}

func (s *SurfaceSession) recordToolPolicy(
	requestID, outcome string,
	call SurfaceToolCall,
	metrics *surfaceRequestMetrics,
) {
	switch outcome {
	case "duplicate_tool_blocked":
		metrics.duplicateToolBlocked++
	case "schema_budget_exhausted":
		metrics.schemaBudgetExhausted++
	case "metadata_tool_failed":
	default:
		outcome = "tool_call_blocked"
	}
	logger.Infof(
		"Agent timing surface=%s request=%s stage=tool_policy tool=%s outcome=%s",
		s.surface.Name(), requestID, call.Name, outcome,
	)
	s.writeAudit("surface_tool_policy", map[string]any{
		"requestId": requestID, "tool": call.Name, "outcome": outcome,
	})
}

func visibleThoughtSummary(value string) string {
	value = strings.TrimSpace(strings.ToValidUTF8(value, "\uFFFD"))
	runes := []rune(value)
	if len(runes) <= maxSurfaceThoughtRunes {
		return value
	}
	return strings.TrimSpace(string(runes[:maxSurfaceThoughtRunes])) + "…"
}

func canonicalSurfaceToolFingerprint(call SurfaceToolCall) string {
	name := strings.ToLower(strings.TrimSpace(call.Name))
	var arguments any
	if err := json.Unmarshal(call.Arguments, &arguments); err != nil {
		return name + "\x00" + strings.TrimSpace(string(call.Arguments))
	}
	canonical, err := json.Marshal(arguments)
	if err != nil {
		return name + "\x00" + strings.TrimSpace(string(call.Arguments))
	}
	return name + "\x00" + string(canonical)
}

func (s *SurfaceSession) complete(
	ctx context.Context,
	request SurfaceRequest,
	state SurfaceState,
	metrics *surfaceRequestMetrics,
) (string, error) {
	tiers := []provider.ContextTier{
		provider.ContextFull, provider.ContextCompact, provider.ContextMinimal,
	}
	var lastErr error
	for index, tier := range tiers {
		if index > 0 {
			s.provider.CompactState(tier)
		}
		callCtx := ctx
		cancel := func() {}
		if timeout := s.config.Provider.RequestTimeout; timeout > 0 {
			callCtx, cancel = context.WithTimeout(ctx, timeout)
		}
		release := func() {}
		queueStarted := time.Now()
		if s.acquireRequest != nil {
			var acquireErr error
			release, acquireErr = s.acquireRequest(callCtx)
			if acquireErr != nil {
				cancel()
				return "", acquireErr
			}
		}
		queueDuration := time.Since(queueStarted)
		metrics.queueDuration += queueDuration
		completionRequest := s.surface.CompletionRequest(request, state, tier)
		metrics.modelRequests++
		modelStarted := time.Now()
		result, err := s.provider.Complete(
			provider.WithLatencyTaskID(callCtx, request.ID),
			completionRequest,
		)
		modelDuration := time.Since(modelStarted)
		metrics.modelDuration += modelDuration
		release()
		cancel()
		outcome := "success"
		if err != nil {
			outcome = "error"
		}
		logger.Infof(
			"Agent timing surface=%s request=%s stage=model round=%d attempt=%d duration_ms=%.3f queue_ms=%.3f outcome=%s",
			s.surface.Name(), request.ID, state.Round, index+1,
			durationMilliseconds(modelDuration), durationMilliseconds(queueDuration), outcome,
		)
		s.emitProgress("Reviewing AI response", "reviewing", "running", true, request.ID, metrics, map[string]any{
			"round": state.Round, "maximumRound": state.MaximumRound,
		})
		if err == nil {
			return result.Content, nil
		}
		lastErr = err
		if !provider.IsKind(err, provider.ErrorContextOverflow) &&
			!provider.IsKind(err, provider.ErrorOutputLimit) {
			return "", err
		}
	}
	return "", lastErr
}

func (s *SurfaceSession) callTool(
	ctx context.Context,
	requestID string,
	state *SurfaceState,
	call SurfaceToolCall,
	metrics *surfaceRequestMetrics,
) error {
	if len(state.ToolResults) >= maxSurfaceToolCalls {
		return fmt.Errorf("agent surface tool call limit reached")
	}
	call.ID = strings.TrimSpace(call.ID)
	if call.ID == "" {
		call.ID = surfaceID("tool")
	}
	if len(call.Arguments) == 0 || len(call.Arguments) > maxSurfaceToolArguments || !json.Valid(call.Arguments) {
		return fmt.Errorf("agent surface tool arguments are invalid")
	}
	if err := s.surface.ValidateTool(call); err != nil {
		return err
	}
	metrics.toolCalls++
	s.emitProgress("Inspecting database metadata", "tool", "running", true, requestID, metrics, map[string]any{
		"tool": call.Name,
	})
	s.writeAudit("surface_tool_request", call)
	started := time.Now()
	result, err := s.tools.Call(ctx, call)
	duration := time.Since(started)
	metrics.toolDuration += duration
	outcome := "success"
	if err != nil {
		outcome = "error"
	}
	logger.Infof(
		"Agent timing surface=%s request=%s stage=tool tool=%s duration_ms=%.3f outcome=%s",
		s.surface.Name(), requestID, call.Name, durationMilliseconds(duration), outcome,
	)
	toolResult := SurfaceToolResult{
		ID: call.ID, Name: call.Name, Arguments: append(json.RawMessage(nil), call.Arguments...),
	}
	if err != nil {
		toolResult.Error = err.Error()
	} else {
		if len(result) > maxSurfaceToolResultBytes || !json.Valid(result) {
			return fmt.Errorf("agent surface tool result is invalid or too large")
		}
		toolResult.Result = append(json.RawMessage(nil), result...)
	}
	state.ToolResults = append(state.ToolResults, toolResult)
	s.writeAudit("surface_tool_result", toolResult)
	if policy, ok := s.surface.(SurfaceToolResultPolicy); ok {
		decision := policy.EvaluateToolResult(*state, toolResult)
		if decision.DisableFurtherTools {
			state.ToolCallsDisabled = true
		}
		if correction := strings.TrimSpace(decision.Correction); correction != "" {
			state.Correction = correction
		}
		if strings.TrimSpace(decision.Outcome) != "" {
			s.recordToolPolicy(requestID, decision.Outcome, call, metrics)
		}
	}
	return nil
}

func (s *SurfaceSession) Interrupt() {
	s.mu.Lock()
	cancel := s.cancel
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (s *SurfaceSession) Close() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	cancel := s.cancel
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	s.lifetimeCancel()
	s.wg.Wait()
	if s.audit != nil {
		s.audit.Close()
	}
}

func (s *SurfaceSession) trimHistoryLocked() {
	total := 0
	index := len(s.history)
	for index > 0 && total < maxSurfaceHistoryBytes {
		index--
		total += len(s.history[index])
	}
	if index > 0 {
		s.history = append(
			[]string{"system: older conversation omitted"}, s.history[index:]...,
		)
	}
}

func (s *SurfaceSession) finishWithError(
	requestID string,
	err error,
	metrics *surfaceRequestMetrics,
) {
	if errors.Is(err, context.Canceled) {
		s.logTiming(requestID, "cancelled", metrics)
		s.emitProgress("Request cancelled", "cancelled", "idle", false, requestID, metrics, nil)
		return
	}
	s.emitData("data-error", map[string]any{
		"code": "agent_failed", "message": err.Error(),
	}, "final", requestID)
	s.logTiming(requestID, "error", metrics)
	s.emitProgress("", "error", "idle", false, requestID, metrics, nil)
	s.writeAudit("surface_error", map[string]any{
		"requestId": requestID, "error": err.Error(),
	})
}

func (s *SurfaceSession) emitProgress(
	text, code, state string,
	interruptible bool,
	requestID string,
	metrics *surfaceRequestMetrics,
	extra map[string]any,
) {
	data := metrics.data()
	data["text"] = text
	data["code"] = code
	data["state"] = state
	data["interruptible"] = interruptible
	for key, value := range extra {
		data[key] = value
	}
	s.emitData("data-progress", data, "process", requestID)
}

func (s *SurfaceSession) emitData(partType string, data any, stage, requestID string) {
	s.emitMessage([]ChatPart{{Type: partType, Data: data}}, stage, requestID)
}

func (s *SurfaceSession) emitMessage(parts []ChatPart, stage, requestID string) {
	message := ChatMessage{
		ID: surfaceID("assistant"), Role: "assistant",
		Metadata: map[string]any{
			"sessionId": s.sessionID, "surface": s.surface.Name(),
			"stage": stage, "requestId": requestID,
		},
		Parts: parts,
	}
	s.emit(message)
}

func (s *SurfaceSession) writeAudit(event string, payload any) {
	if s.audit != nil {
		s.audit.Write(event, payload)
	}
}

func (s *SurfaceSession) logTiming(
	requestID, outcome string,
	metrics *surfaceRequestMetrics,
) {
	logger.Infof(
		"Agent timing surface=%s request=%s stage=request duration_ms=%.3f model_requests=%d model_ms=%.3f tool_calls=%d tool_ms=%.3f queue_ms=%.3f duplicate_tool_blocked=%d schema_budget_exhausted=%d forced_clarifications=%d outcome=%s",
		s.surface.Name(), requestID, durationMilliseconds(time.Since(metrics.started)),
		metrics.modelRequests, durationMilliseconds(metrics.modelDuration), metrics.toolCalls,
		durationMilliseconds(metrics.toolDuration), durationMilliseconds(metrics.queueDuration),
		metrics.duplicateToolBlocked, metrics.schemaBudgetExhausted, metrics.forcedClarifications, outcome,
	)
}

var surfaceSequence atomic.Uint64

func surfaceID(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixMilli(), surfaceSequence.Add(1))
}
