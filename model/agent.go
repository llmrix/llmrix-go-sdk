package model

// Agent represents a Llmrix agent definition (cloud mode only).
type Agent struct {
	ID           int                      `json:"id"`
	Name         string                   `json:"name"`
	Type         string                   `json:"type"` // "solo" | "team"
	SystemPrompt string                   `json:"system_prompt"`
	ModelID      string                   `json:"model_id"`
	MCPs         []string                 `json:"mcps"`
	Skills       []string                 `json:"skills"`
	Metadata     map[string]interface{}   `json:"metadata"`
	Introduce    string                   `json:"introduce"`
	Outlines     []map[string]interface{} `json:"outlines"`
}

// AgentCreateRequest is the request body for creating an agent.
type AgentCreateRequest struct {
	Name         string                   `json:"name"`
	Type         string                   `json:"type,omitempty"`
	SystemPrompt string                   `json:"system_prompt,omitempty"`
	ModelID      string                   `json:"model_id,omitempty"`
	MCPs         []string                 `json:"mcps,omitempty"`
	Skills       []string                 `json:"skills,omitempty"`
	Metadata     map[string]interface{}   `json:"metadata,omitempty"`
	Introduce    string                   `json:"introduce,omitempty"`
	Outlines     []map[string]interface{} `json:"outlines,omitempty"`
}

// AgentUpdateRequest is the request body for updating an agent.
// Only non-nil pointer fields are sent to the server.
type AgentUpdateRequest struct {
	Name         *string                  `json:"name,omitempty"`
	Type         *string                  `json:"type,omitempty"`
	SystemPrompt *string                  `json:"system_prompt,omitempty"`
	ModelID      *string                  `json:"model_id,omitempty"`
	MCPs         []string                 `json:"mcps,omitempty"`
	Skills       []string                 `json:"skills,omitempty"`
	Metadata     map[string]interface{}   `json:"metadata,omitempty"`
	Introduce    *string                  `json:"introduce,omitempty"`
	Outlines     []map[string]interface{} `json:"outlines,omitempty"`
}

// Mate represents a sub-agent assigned to a team agent.
type Mate struct {
	Name         string `json:"name"`
	Introduce    string `json:"introduce,omitempty"`
	SystemPrompt string `json:"system_prompt,omitempty"`
	Sort         int    `json:"sort,omitempty"`
}

// SaveMatesRequest is the wire body for the save-mates endpoint.
type SaveMatesRequest struct {
	Mates []Mate `json:"mates"`
}
