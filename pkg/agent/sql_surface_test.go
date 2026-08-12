package agent

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSQLSurfaceProposalRequiresValidationAndPreservesEditorBase(t *testing.T) {
	surface := NewSQLSurface()
	contextJSON := json.RawMessage(`{
        "dialect":"postgresql","database":"app","schema":"public","nodeKey":"node-1",
        "paneId":"pane-1","tabId":"tab-1","workspaceTabId":"tab-1","workspaceTabKind":"query",
        "currentContext":"public","revision":7,
        "selectionFrom":7,"selectionTo":19,"selectedSql":"* FROM users",
        "documentSql":"SELECT * FROM users"
    }`)
	request := SurfaceRequest{ID: "request-1", Operation: "repair", Context: contextJSON}
	action, err := surface.DecodeAction(`{
        "kind":"proposal","message":"Use explicit columns","thoughtSummary":"Validated SQL",
		"toolName":"","toolArguments":{"query":"","schema":"","tables":[],"sql":""},
        "sql":"id, name FROM users","proposalExplanation":"Avoid SELECT star",
        "analysis":{"valid":true,"statementCount":1,"statementType":"SELECT","riskLevel":1,
        "riskReason":"Read-only SQL statement","tables":["users"],"columns":["id","name"],"errors":[]}
    }`)
	if err != nil {
		t.Fatal(err)
	}

	review, err := surface.Review(request, SurfaceState{}, action)
	if err != nil {
		t.Fatal(err)
	}
	if review.Tool == nil || review.Tool.Name != "validate_sql" || !review.FinalizeAfterTool {
		t.Fatalf("review tool = %#v, want validate_sql", review.Tool)
	}

	state := SurfaceState{ToolResults: []SurfaceToolResult{{
		Name:      "validate_sql",
		Arguments: json.RawMessage(`{"sql":"id, name FROM users"}`),
		Result: json.RawMessage(`{
            "valid":true,"statementCount":1,"statementType":"SELECT","riskLevel":1,
            "riskReason":"Read-only SQL statement","tables":["users"],"columns":["id","name"],"errors":[]
        }`),
	}}}
	review, err = surface.Review(request, state, action)
	if err != nil || review.Tool != nil || review.Correction != "" {
		t.Fatalf("validated review = %#v, err=%v", review, err)
	}

	parts, err := surface.FinalParts(request, state, action)
	if err != nil {
		t.Fatal(err)
	}
	var proposal map[string]any
	for _, part := range parts {
		if part.Type == "data-sql-proposal" {
			proposal = part.Data.(map[string]any)
		}
	}
	if proposal == nil || proposal["originalSql"] != "* FROM users" {
		t.Fatalf("proposal = %#v", proposal)
	}
	if proposal["sql"] != "id, name FROM users" {
		t.Fatalf("proposal sql = %#v, want selection replacement only", proposal["sql"])
	}
	base := proposal["base"].(map[string]any)
	if base["target"] != "selection" || base["paneId"] != "pane-1" || base["tabId"] != "tab-1" || base["revision"] != int64(7) {
		t.Fatalf("proposal base = %#v", base)
	}
	if base["workspaceTabId"] != "tab-1" || base["workspaceTabKind"] != "query" ||
		base["currentContext"] != "public" {
		t.Fatalf("proposal workspace base = %#v", base)
	}
}

func TestSQLSurfaceOmitsUnchangedProposal(t *testing.T) {
	surface := NewSQLSurface()
	request := SurfaceRequest{Operation: "repair", Context: json.RawMessage(`{
        "nodeKey":"node-1","paneId":"pane-1","tabId":"tab-1",
        "selectionFrom":7,"selectionTo":19,"selectedSql":"* FROM users",
        "documentSql":"SELECT * FROM users"
    }`)}
	action, err := surface.DecodeAction(`{
        "kind":"proposal","message":"The selected SQL is already valid","thoughtSummary":"Validated SQL",
		"toolName":"","toolArguments":{"query":"","schema":"","tables":[],"sql":""},
        "sql":"* FROM users","proposalExplanation":"No repair is needed",
        "analysis":{"valid":true,"statementCount":1,"statementType":"SELECT","riskLevel":1,
        "riskReason":"Read-only SQL statement","tables":["users"],"columns":[],"errors":[]}
    }`)
	if err != nil {
		t.Fatal(err)
	}
	state := SurfaceState{ToolResults: []SurfaceToolResult{{
		Name:      "validate_sql",
		Arguments: json.RawMessage(`{"sql":"* FROM users"}`),
		Result:    json.RawMessage(`{"valid":true,"statementCount":1,"statementType":"SELECT","riskLevel":1,"errors":[]}`),
	}}}
	parts, err := surface.FinalParts(request, state, action)
	if err != nil {
		t.Fatal(err)
	}
	var hasText, hasAnalysis bool
	for _, part := range parts {
		hasText = hasText || part.Type == "text"
		hasAnalysis = hasAnalysis || part.Type == "data-sql-analysis"
		if part.Type == "data-sql-proposal" {
			t.Fatalf("unchanged SQL emitted proposal: %#v", part.Data)
		}
	}
	if !hasText || !hasAnalysis {
		t.Fatalf("unchanged SQL parts = %#v, want explanation and analysis", parts)
	}
}

func TestSQLSurfaceRejectsMultiStatementProposal(t *testing.T) {
	surface := NewSQLSurface()
	request := SurfaceRequest{Context: json.RawMessage(`{"documentSql":"SELECT 1"}`)}
	action, err := surface.DecodeAction(`{
        "kind":"proposal","message":"","thoughtSummary":"","toolName":"",
		"toolArguments":{"query":"","schema":"","tables":[],"sql":""},
        "sql":"SELECT 1; SELECT 2","proposalExplanation":"",
        "analysis":{"valid":true,"statementCount":2,"statementType":"MULTI","riskLevel":1,
        "riskReason":"","tables":[],"columns":[],"errors":[]}
    }`)
	if err != nil {
		t.Fatal(err)
	}
	state := SurfaceState{ToolResults: []SurfaceToolResult{{
		Name:      "validate_sql",
		Arguments: json.RawMessage(`{"sql":"SELECT 1; SELECT 2"}`),
		Result:    json.RawMessage(`{"valid":true,"statementCount":2,"errors":[]}`),
	}}}
	review, err := surface.Review(request, state, action)
	if err != nil {
		t.Fatal(err)
	}
	if review.Correction == "" {
		t.Fatal("multi-statement proposal was not rejected")
	}
}

func TestSQLSurfaceAllowsOnlyMetadataAndValidationTools(t *testing.T) {
	surface := NewSQLSurface()
	if err := surface.ValidateTool(SurfaceToolCall{Name: "inspect_schema", Arguments: json.RawMessage(`{"tables":["users"]}`)}); err != nil {
		t.Fatal(err)
	}
	if err := surface.ValidateTool(SurfaceToolCall{Name: "describe_table", Arguments: json.RawMessage(`{"table":"users"}`)}); err == nil {
		t.Fatal("legacy single-table metadata tool should not be exposed")
	}
	if err := surface.ValidateTool(SurfaceToolCall{Name: "query_rows", Arguments: json.RawMessage(`{}`)}); err == nil {
		t.Fatal("business-row tool should not be allowed")
	}
}

func TestSQLSurfaceModelOnlySeesBatchedMetadataTool(t *testing.T) {
	tool := sqlAssistantActionTool()
	properties := tool.Parameters["properties"].(map[string]any)
	toolName := properties["toolName"].(map[string]any)
	names := toolName["enum"].([]string)
	if len(names) != 2 || names[0] != "" || names[1] != "inspect_schema" {
		t.Fatalf("model tool names = %#v, want only batched metadata inspection", names)
	}
}

func TestSQLSurfaceDecodeRejectsIncompleteToolAction(t *testing.T) {
	surface := NewSQLSurface()
	_, err := surface.DecodeAction(`{
		"kind":"tool","message":"","thoughtSummary":"Searching schema","toolName":"inspect_schema",
		"toolArguments":{"query":"","schema":"","tables":[],"sql":""},
		"sql":"","proposalExplanation":"",
		"analysis":{"valid":false,"statementType":"","riskLevel":0,
		"riskReason":"","tables":[],"columns":[],"errors":[]}
	}`)
	if err == nil || !strings.Contains(err.Error(), "schema inspection query or tables are required") {
		t.Fatalf("DecodeAction error = %v, want missing schema inspection target", err)
	}
}

func TestSQLSurfacePrefetchesReferencedTablesInOneToolCall(t *testing.T) {
	surface := NewSQLSurface()
	calls, err := surface.InitialTools(SurfaceRequest{Operation: "generate", Context: json.RawMessage(`{
		"schema":"public","referencedTables":["users","orders","USERS",""]
	}`)})
	if err != nil {
		t.Fatal(err)
	}
	if len(calls) != 1 || calls[0].Name != "inspect_schema" {
		t.Fatalf("initial calls = %#v", calls)
	}
	var arguments sqlToolArguments
	if err = json.Unmarshal(calls[0].Arguments, &arguments); err != nil {
		t.Fatal(err)
	}
	if arguments.Schema != "public" || len(arguments.Tables) != 2 ||
		arguments.Tables[0] != "users" || arguments.Tables[1] != "orders" {
		t.Fatalf("inspect arguments = %#v", arguments)
	}
}

func TestSQLSurfaceDisablesReasoningAfterFirstModelRound(t *testing.T) {
	surface := NewSQLSurface()
	request := SurfaceRequest{Operation: "generate", Context: json.RawMessage(`{}`)}
	first := surface.CompletionRequest(request, SurfaceState{Round: 1}, "full")
	second := surface.CompletionRequest(request, SurfaceState{Round: 2}, "full")
	if first.ReasoningMode != "" {
		t.Fatalf("first round reasoning mode = %q, want provider default", first.ReasoningMode)
	}
	if second.ReasoningMode != "off" {
		t.Fatalf("second round reasoning mode = %q, want off", second.ReasoningMode)
	}
}
