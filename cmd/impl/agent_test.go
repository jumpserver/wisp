package impl

import (
	"context"
	"strings"
	"testing"
	"time"

	protobuf "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func TestAgentRequestLimiterUsesBoundedWaitQueue(t *testing.T) {
	limiter := newAgentRequestLimiter(1, 1)
	releaseFirst, err := limiter.acquire(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	waiter := make(chan error, 1)
	waiterRelease := make(chan func(), 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		release, acquireErr := limiter.acquire(ctx)
		if acquireErr == nil {
			waiterRelease <- release
		}
		waiter <- acquireErr
	}()

	deadline := time.Now().Add(time.Second)
	for limiter.waiting.Load() != 1 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if limiter.waiting.Load() != 1 {
		t.Fatal("request did not enter the wait queue")
	}

	_, err = limiter.acquire(context.Background())
	if err == nil || !strings.Contains(err.Error(), "queue is full") {
		t.Fatalf("overflow error = %v", err)
	}

	releaseFirst()
	if err = <-waiter; err != nil {
		t.Fatal(err)
	}
	(<-waiterRelease)()
}

func TestAgentRequestLimiterIsUnlimitedByDefault(t *testing.T) {
	if limiter := newAgentRequestLimiter(0, 100); limiter != nil {
		t.Fatalf("limiter = %#v, want nil", limiter)
	}
}

func TestAgentRequestLimiterCanDisableWaiting(t *testing.T) {
	limiter := newAgentRequestLimiter(1, 0)
	release, err := limiter.acquire(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if _, err = limiter.acquire(context.Background()); err == nil {
		t.Fatal("request should be rejected when the wait queue is disabled")
	}
}

func TestAgentToolBridgeIgnoresLateCancelledResult(t *testing.T) {
	bridge := &agentToolBridge{pending: make(map[string]chan agentToolReply)}
	if err := bridge.resolve(&protobuf.AgentToolResult{Id: "cancelled", ResultJson: `{}`}); err != nil {
		t.Fatal(err)
	}
}
