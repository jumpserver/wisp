package agent

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jumpserver-dev/sdk-go/model"
	"github.com/jumpserver/wisp/pkg/agent/provider"
)

const (
	providerEnvName           = "TERMINAL_AI_PROVIDER"
	toolCallEnvName           = "TERMINAL_AI_TOOL_CALL"
	maxAgentModelOutputTokens = int64(32 * 1024)
)

type Config struct {
	Provider         provider.Config
	MemoryRoot       string
	MemorySessions   int
	MaxModelRequests int
}

func NewConfig(modelConfig model.TerminalConfig) Config {
	return NewConfigWithDataRoot(modelConfig, "", false)
}

func NewConfigWithDataRoot(
	modelConfig model.TerminalConfig,
	dataRoot string,
	auditEnabled bool,
) Config {
	name := strings.TrimSpace(modelConfig.ChatAIProvider)
	if name == "" {
		name = strings.TrimSpace(os.Getenv(providerEnvName))
	}
	if name == "" {
		name = provider.NameOpenAICompatible
	}
	providerConfig := provider.NormalizeConfig(provider.Config{
		Name: name, APIKey: modelConfig.ChatAIApiKey,
		BaseURL: modelConfig.ChatAIBaseUrl, Model: modelConfig.ChatAIModel,
		Proxy: modelConfig.ChatAIProxy, ToolCallMode: os.Getenv(toolCallEnvName),
		ReasoningMode: provider.ReasoningAuto,
		Store:         false, NativeCompaction: false,
		ContextSoftLimitPercent: 80, RequestTimeout: 5 * time.Minute,
	})
	if providerConfig.MaxOutputTokens > maxAgentModelOutputTokens {
		providerConfig.MaxOutputTokens = maxAgentModelOutputTokens
	}
	config := Config{
		Provider:       providerConfig,
		MemorySessions: 10, MaxModelRequests: 30,
	}
	if auditEnabled && strings.TrimSpace(dataRoot) != "" {
		config.MemoryRoot = filepath.Join(dataRoot, "agent", "audit")
	}
	return config
}
