package protobuf

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestAgentSessionProtocolIsBidirectionalAndRoundTripsIdentity(t *testing.T) {
	var agentStreamFound bool
	for _, stream := range Service_ServiceDesc.Streams {
		if stream.StreamName != "AgentSession" {
			continue
		}
		agentStreamFound = true
		if !stream.ClientStreams || !stream.ServerStreams {
			t.Fatalf("AgentSession stream descriptor = %#v", stream)
		}
	}
	if !agentStreamFound {
		t.Fatal("AgentSession stream descriptor is missing")
	}

	original := &AgentClientEvent{Event: &AgentClientEvent_Open{Open: &AgentSessionOpen{
		SessionId: "session-1", UserId: "user-1", OrganizationId: "org-1",
		AssetId: "asset-1", AccountId: "account-1", Protocol: "postgresql",
		Language: "zh-CN", Surface: "sql",
	}}}
	encoded, err := proto.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	decoded := new(AgentClientEvent)
	if err = proto.Unmarshal(encoded, decoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(original, decoded) {
		t.Fatalf("decoded event = %#v", decoded)
	}
}
