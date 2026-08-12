package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/jumpserver/wisp/pkg/agent"
	appconfig "github.com/jumpserver/wisp/pkg/config"
	"github.com/jumpserver/wisp/pkg/logger"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

type agentRequestLimiter struct {
	tokens   chan struct{}
	maxQueue int64
	waiting  atomic.Int64
}

func newAgentRequestLimiter(maxConcurrent, maxQueue int) *agentRequestLimiter {
	if maxConcurrent <= 0 {
		return nil
	}
	if maxQueue < 0 {
		maxQueue = 0
	}
	return &agentRequestLimiter{
		tokens: make(chan struct{}, maxConcurrent), maxQueue: int64(maxQueue),
	}
}

func (l *agentRequestLimiter) acquire(ctx context.Context) (func(), error) {
	if l == nil {
		return func() {}, nil
	}
	select {
	case l.tokens <- struct{}{}:
		return func() { <-l.tokens }, nil
	default:
	}
	waiting := l.waiting.Add(1)
	if waiting > l.maxQueue {
		l.waiting.Add(-1)
		return nil, fmt.Errorf("agent request queue is full")
	}
	defer l.waiting.Add(-1)
	select {
	case l.tokens <- struct{}{}:
		return func() { <-l.tokens }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type agentSessionRegistry struct {
	sync.Mutex
	items map[string]*agent.SurfaceSession
}

func newAgentSessionRegistry() *agentSessionRegistry {
	return &agentSessionRegistry{items: make(map[string]*agent.SurfaceSession)}
}

func (r *agentSessionRegistry) store(id string, session *agent.SurfaceSession) bool {
	r.Lock()
	defer r.Unlock()
	if _, exists := r.items[id]; exists {
		return false
	}
	r.items[id] = session
	return true
}

func (r *agentSessionRegistry) remove(id string, session *agent.SurfaceSession) {
	r.Lock()
	if r.items[id] == session {
		delete(r.items, id)
	}
	r.Unlock()
}

func (r *agentSessionRegistry) close(id string) {
	r.Lock()
	session := r.items[id]
	delete(r.items, id)
	r.Unlock()
	if session != nil {
		session.Close()
	}
}

type agentStreamSender struct {
	stream pb.Service_AgentSessionServer
	sync.Mutex
}

func (s *agentStreamSender) send(event *pb.AgentServerEvent) error {
	s.Lock()
	defer s.Unlock()
	return s.stream.Send(event)
}

func (s *agentStreamSender) ready(value *pb.AgentReady) error {
	return s.send(&pb.AgentServerEvent{
		Event: &pb.AgentServerEvent_Ready{Ready: value},
	})
}

func (s *agentStreamSender) failure(code, message, requestID string) error {
	return s.send(&pb.AgentServerEvent{
		Event: &pb.AgentServerEvent_Error{Error: &pb.AgentError{
			Code: code, Message: message, RequestId: requestID,
		}},
	})
}

type agentToolReply struct {
	result json.RawMessage
	err    error
}

type agentToolBridge struct {
	sender *agentStreamSender

	sync.Mutex
	pending map[string]chan agentToolReply
	closed  bool
}

func newAgentToolBridge(sender *agentStreamSender) *agentToolBridge {
	return &agentToolBridge{
		sender: sender, pending: make(map[string]chan agentToolReply),
	}
}

func (b *agentToolBridge) Call(
	ctx context.Context,
	call agent.SurfaceToolCall,
) (json.RawMessage, error) {
	reply := make(chan agentToolReply, 1)
	b.Lock()
	if b.closed {
		b.Unlock()
		return nil, fmt.Errorf("agent tool bridge is closed")
	}
	if _, exists := b.pending[call.ID]; exists {
		b.Unlock()
		return nil, fmt.Errorf("duplicate agent tool call %s", call.ID)
	}
	b.pending[call.ID] = reply
	b.Unlock()

	defer func() {
		b.Lock()
		delete(b.pending, call.ID)
		b.Unlock()
	}()
	if err := b.sender.send(&pb.AgentServerEvent{
		Event: &pb.AgentServerEvent_ToolCall{ToolCall: &pb.AgentToolCall{
			Id: call.ID, Name: call.Name, ArgumentsJson: string(call.Arguments),
		}},
	}); err != nil {
		return nil, err
	}
	select {
	case value := <-reply:
		return value.result, value.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (b *agentToolBridge) resolve(value *pb.AgentToolResult) error {
	if value == nil || strings.TrimSpace(value.Id) == "" {
		return fmt.Errorf("agent tool result id is required")
	}
	b.Lock()
	reply := b.pending[value.Id]
	b.Unlock()
	if reply == nil {
		// A cancelled request may finish its bounded JDBC metadata call after
		// the runtime has removed the correlation. The late result cannot affect
		// any request and is safe to discard.
		return nil
	}
	result := agentToolReply{}
	if value.Error != "" {
		result.err = fmt.Errorf("Chen metadata tool failed: %s", value.Error)
	} else {
		result.result = json.RawMessage(value.ResultJson)
		if len(result.result) == 0 || !json.Valid(result.result) {
			result.err = fmt.Errorf("Chen metadata tool returned invalid JSON")
		}
	}
	select {
	case reply <- result:
		return nil
	default:
		return nil
	}
}

func (b *agentToolBridge) close() {
	b.Lock()
	b.closed = true
	for id, reply := range b.pending {
		select {
		case reply <- agentToolReply{err: fmt.Errorf("agent tool bridge closed")}:
		default:
		}
		delete(b.pending, id)
	}
	b.Unlock()
}

func (j *JMServer) AgentSession(stream pb.Service_AgentSessionServer) error {
	sender := &agentStreamSender{stream: stream}
	first, err := stream.Recv()
	if err != nil {
		return err
	}
	open := first.GetOpen()
	if open == nil {
		_ = sender.failure("protocol_error", "the first agent event must open a session", "")
		return nil
	}
	if err = j.validateAgentSession(open); err != nil {
		_ = sender.ready(&pb.AgentReady{
			Enabled: false, Reason: err.Error(), SessionId: open.SessionId,
			Surface: open.Surface,
		})
		return nil
	}

	termConfig := j.uploader.GetTerminalSetting()
	localConfig := appconfig.Get()
	agentConfig := agent.NewConfigWithDataRoot(
		termConfig, localConfig.DataFolderPath, localConfig.AIAuditEnabled,
	)
	bridge := newAgentToolBridge(sender)
	var surface agent.Surface
	switch strings.ToLower(strings.TrimSpace(open.Surface)) {
	case agent.SQLSurfaceName:
		surface = agent.NewSQLSurface()
	default:
		_ = sender.ready(&pb.AgentReady{
			Enabled: false, Reason: "unsupported agent surface", SessionId: open.SessionId,
			Surface: open.Surface,
		})
		return nil
	}

	session, err := agent.NewSurfaceSession(agent.SurfaceSessionOptions{
		SessionID:      open.SessionId,
		UserID:         open.UserId,
		Language:       open.Language,
		Config:         agentConfig,
		Surface:        surface,
		Tools:          bridge,
		AcquireRequest: j.agentLimiter.acquire,
		Emit: func(message agent.ChatMessage) {
			value, marshalErr := json.Marshal(message)
			if marshalErr != nil {
				return
			}
			if sendErr := sender.send(&pb.AgentServerEvent{
				Event: &pb.AgentServerEvent_Chat{Chat: &pb.AgentChatMessage{
					MessageJson: string(value),
				}},
			}); sendErr != nil {
				logger.Errorf("Send agent chat event failed for session %s: %s", open.SessionId, sendErr)
			}
		},
	})
	if err != nil {
		_ = sender.ready(&pb.AgentReady{
			Enabled: false, Reason: err.Error(), SessionId: open.SessionId,
			Surface: open.Surface,
		})
		return nil
	}
	if !j.agentSessions.store(open.SessionId, session) {
		session.Close()
		_ = sender.ready(&pb.AgentReady{
			Enabled: false, Reason: "an agent session is already active", SessionId: open.SessionId,
			Surface: open.Surface,
		})
		return nil
	}
	defer func() {
		j.agentSessions.remove(open.SessionId, session)
		session.Close()
		bridge.close()
	}()

	info := session.ProviderInfo()
	if err = sender.ready(&pb.AgentReady{
		Enabled: true, SessionId: open.SessionId, Surface: surface.Name(),
		Provider: info.Name, Model: info.Model,
	}); err != nil {
		return err
	}
	session.AnnounceCapability()

	for {
		event, recvErr := stream.Recv()
		if recvErr == io.EOF {
			return nil
		}
		if recvErr != nil {
			return recvErr
		}
		switch {
		case event.GetRequest() != nil:
			request := event.GetRequest()
			handleErr := session.Handle(agent.SurfaceRequest{
				ID: request.Id, Operation: request.Operation,
				Question: request.Question, Context: json.RawMessage(request.ContextJson),
			})
			if handleErr != nil {
				_ = sender.failure("invalid_request", handleErr.Error(), request.Id)
			}
		case event.GetToolResult() != nil:
			if resolveErr := bridge.resolve(event.GetToolResult()); resolveErr != nil {
				_ = sender.failure("invalid_tool_result", resolveErr.Error(), "")
			}
		case event.GetCancel() != nil:
			session.Interrupt()
		case event.GetClose() != nil:
			return nil
		default:
			_ = sender.failure("protocol_error", "unsupported agent client event", "")
		}
	}
}

func (j *JMServer) validateAgentSession(open *pb.AgentSessionOpen) error {
	if open == nil || strings.TrimSpace(open.SessionId) == "" {
		return fmt.Errorf("agent session id is required")
	}
	session, ok := j.beat.GetSession(open.SessionId)
	if !ok {
		return fmt.Errorf("the associated JMS session is not active")
	}
	checks := []struct {
		name     string
		expected string
		actual   string
	}{
		{"user", session.UserID, open.UserId},
		{"organization", session.OrgID, open.OrganizationId},
		{"asset", session.AssetID, open.AssetId},
		{"account", session.AccountID, open.AccountId},
	}
	for _, check := range checks {
		if strings.TrimSpace(check.expected) == "" || check.expected != check.actual {
			return fmt.Errorf("agent session %s does not match the JMS session", check.name)
		}
	}
	if !strings.EqualFold(session.Protocol, open.Protocol) {
		return fmt.Errorf("agent session protocol does not match the JMS session")
	}
	if strings.EqualFold(open.Surface, agent.SQLSurfaceName) {
		switch strings.ToLower(strings.TrimSpace(open.Protocol)) {
		case agent.ProtocolMySQL, agent.ProtocolMariaDB, agent.ProtocolPostgreSQL,
			agent.ProtocolSQLServer, agent.ProtocolOracle, agent.ProtocolClickHouse,
			agent.ProtocolDameng, agent.ProtocolDB2:
		default:
			return fmt.Errorf("SQL agent surface does not support protocol %s", open.Protocol)
		}
	}
	return nil
}
