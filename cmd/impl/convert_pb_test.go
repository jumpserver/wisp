package impl

import (
	"testing"

	"github.com/jumpserver-dev/sdk-go/model"
)

func TestConvertToProtobufSessionPreservesIdentity(t *testing.T) {
	session := model.Session{
		ID:        "session-id",
		OrgID:     "organization-id",
		UserID:    "user-id",
		AssetID:   "asset-id",
		AccountID: "account-id",
		TokenId:   "token-id",
	}

	converted := ConvertToProtobufSession(session)
	if converted.GetAccountId() != session.AccountID {
		t.Fatalf("account id = %q, want %q", converted.GetAccountId(), session.AccountID)
	}
	if converted.GetTokenId() != session.TokenId {
		t.Fatalf("token id = %q, want %q", converted.GetTokenId(), session.TokenId)
	}
}

func TestConvertToPbSettingPreservesChatAIGate(t *testing.T) {
	converted := ConvertToPbSetting(&model.TerminalConfig{}, true)
	if !converted.GetChatAiEnabled() {
		t.Fatal("chat AI feature gate was not preserved")
	}
}

func TestChatAIEnabledByModelConfig(t *testing.T) {
	tests := []struct {
		name    string
		setting model.TerminalConfig
		enabled bool
	}{
		{name: "missing configuration"},
		{
			name: "disabled with complete configuration",
			setting: model.TerminalConfig{
				ChatAIApiKey: "key", ChatAIModel: "model",
			},
		},
		{
			name: "missing model",
			setting: model.TerminalConfig{
				ChatAIEnabled: true, ChatAIApiKey: "key",
			},
		},
		{
			name: "missing api key",
			setting: model.TerminalConfig{
				ChatAIEnabled: true, ChatAIModel: "model",
			},
		},
		{
			name: "configured without base url",
			setting: model.TerminalConfig{
				ChatAIEnabled: true,
				ChatAIApiKey:  " key ", ChatAIModel: " model ",
			},
			enabled: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := chatAIEnabledByModelConfig(test.setting); got != test.enabled {
				t.Fatalf("enabled = %t, want %t", got, test.enabled)
			}
		})
	}
}
