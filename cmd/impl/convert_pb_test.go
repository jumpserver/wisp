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
