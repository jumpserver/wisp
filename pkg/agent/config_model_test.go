package agent

import (
	"testing"

	"github.com/jumpserver-dev/sdk-go/model"
	"github.com/jumpserver/wisp/pkg/agent/provider"
)

func TestConfigPrefersChatAIProvider(t *testing.T) {
	t.Setenv(providerEnvName, provider.NameOpenAI)
	config := NewConfig(model.TerminalConfig{ChatAIProvider: provider.NameDeepSeek})
	if config.Provider.Name != provider.NameDeepSeek {
		t.Fatalf("provider = %q, want ChatAIProvider", config.Provider.Name)
	}
	config = NewConfig(model.TerminalConfig{})
	if config.Provider.Name != provider.NameOpenAI {
		t.Fatalf("provider = %q, want environment fallback", config.Provider.Name)
	}
	if config.Provider.Store || config.Provider.ReasoningMode != provider.ReasoningAuto ||
		config.MemorySessions != 10 {
		t.Fatalf("unexpected Terminal AI defaults: %#v", config)
	}
}
