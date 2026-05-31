// Package model contains all domain types used by the Llmrix SDK.
// Types are re-exported from the root llmrix package for convenience.
package model

// Conversation represents a Llmrix chat thread.
type Conversation struct {
	// Seq is the server-assigned cursor value for keyset pagination.
	Seq string `json:"seq"`
	// ID is the unique conversation / thread identifier.
	ID string `json:"id"`
	// Title is the human-readable conversation title.
	Title string `json:"title"`
	// AgentID is the ID of the bound agent; empty if none.
	AgentID string `json:"agent_id"`
	// CreatedAt is the ISO-8601 creation timestamp.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the ISO-8601 last-update timestamp.
	UpdatedAt string `json:"updated_at"`
}

// ConversationCreateRequest is the request body for creating a conversation.
type ConversationCreateRequest struct {
	Title   string `json:"title"`
	AgentID string `json:"agent_id,omitempty"`
}

// ConversationUpdateRequest is the request body for updating a conversation.
type ConversationUpdateRequest struct {
	Title string `json:"title,omitempty"`
}

// Message is a single message within a conversation.
type Message struct {
	Seq       int        `json:"seq"`
	ID        string     `json:"id"`
	Role      string     `json:"role"` // "user" | "assistant" | "tool" | "system"
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall is a tool invocation embedded in an assistant message.
type ToolCall struct {
	ID   string                 `json:"id"`
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

// PageResult is the generic paged list returned by list endpoints.
type PageResult[T any] struct {
	Items   []T  `json:"items"`
	HasMore bool `json:"has_more"`
}
