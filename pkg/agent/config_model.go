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
	name := strings.TrimSpace(modelConfig.ChatAIType)
	if name == "" {
		name = strings.TrimSpace(os.Getenv(providerEnvName))
	}
	if name == "" {
		name = provider.NameGPT
	}
	providerConfig := provider.NormalizeConfig(provider.Config{
		Name: name, APIKey: modelConfig.GptApiKey,
		BaseURL: modelConfig.GptBaseUrl, Model: modelConfig.GptModel,
		Proxy: modelConfig.GptProxy, ToolCallMode: os.Getenv(toolCallEnvName),
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
