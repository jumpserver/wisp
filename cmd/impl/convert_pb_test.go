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

func TestChatAIEnabledByURL(t *testing.T) {
	if chatAIEnabledByURL(model.TerminalConfig{GptBaseUrl: " \t "}) {
		t.Fatal("chat AI must remain disabled when GPT_BASE_URL is empty")
	}
	if !chatAIEnabledByURL(model.TerminalConfig{GptBaseUrl: "  https://ai.example.test/v1  "}) {
		t.Fatal("chat AI must be enabled when GPT_BASE_URL is configured")
	}
}
