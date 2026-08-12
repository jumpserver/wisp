package agent

type ChatMessage struct {
	ID       string         `json:"id"`
	Role     string         `json:"role"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Parts    []ChatPart     `json:"parts"`
}

type ChatPart struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	State string `json:"state,omitempty"`
	Data  any    `json:"data,omitempty"`
}
