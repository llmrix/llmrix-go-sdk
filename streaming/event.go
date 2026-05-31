// Package streaming defines all SSE event types emitted during an agent chat run.
//
// Use a type switch to handle specific events:
//
//	switch e := event.(type) {
//	case *streaming.MessageChunkEvent:
//	    fmt.Print(e.Content)
//	case *streaming.RunEndEvent:
//	    fmt.Println("\n[done]")
//	}
//
// All concrete types are re-exported from the root llmrix package for convenience.
package streaming

// StreamEvent is the sealed interface implemented by every SSE event type.
type StreamEvent interface {
	// Channel returns the SSE channel this event belongs to.
	Channel() string
	// Type returns the SSE type discriminator, or "" for channel-only events.
	Type() string
}

// ---------------------------------------------------------------------------
// Lifecycle
// ---------------------------------------------------------------------------

// RunStartEvent is emitted at the beginning of every agent turn.
type RunStartEvent struct {
	RunID    string `json:"run_id"`
	ThreadID string `json:"thread_id"`
}

func (*RunStartEvent) Channel() string { return "lifecycle" }
func (*RunStartEvent) Type() string    { return "run_start" }

// RunEndEvent is emitted when the agent run completes normally.
type RunEndEvent struct {
	RunID    string `json:"run_id"`
	ThreadID string `json:"thread_id"`
}

func (*RunEndEvent) Channel() string { return "lifecycle" }
func (*RunEndEvent) Type() string    { return "run_end" }

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

// MessageChunkEvent carries an incremental text fragment of the assistant reply.
// Concatenate successive Content values to reconstruct the full message.
type MessageChunkEvent struct {
	Content string `json:"content"`
	ID      string `json:"id"`
}

func (*MessageChunkEvent) Channel() string { return "messages" }
func (*MessageChunkEvent) Type() string    { return "message_chunk" }

// ---------------------------------------------------------------------------
// Tools
// ---------------------------------------------------------------------------

// ToolStartEvent is emitted when the agent begins executing a tool call.
type ToolStartEvent struct {
	ID   string                 `json:"id"`
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

func (*ToolStartEvent) Channel() string { return "tools" }
func (*ToolStartEvent) Type() string    { return "tool_start" }

// ToolEndEvent is emitted when a tool call finishes.
type ToolEndEvent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
	Error  bool   `json:"error"`
}

func (*ToolEndEvent) Channel() string { return "tools" }
func (*ToolEndEvent) Type() string    { return "tool_end" }

// SubagentStartEvent is emitted when the agent spawns a sub-agent.
type SubagentStartEvent struct {
	ID        string `json:"id"`
	AgentName string `json:"agent_name"`
}

func (*SubagentStartEvent) Channel() string { return "tools" }
func (*SubagentStartEvent) Type() string    { return "subagent_start" }

// SubagentEndEvent is emitted when a sub-agent completes its task.
type SubagentEndEvent struct {
	ID        string `json:"id"`
	AgentName string `json:"agent_name"`
	Error     bool   `json:"error"`
}

func (*SubagentEndEvent) Channel() string { return "tools" }
func (*SubagentEndEvent) Type() string    { return "subagent_end" }

// ---------------------------------------------------------------------------
// HITL (Human-In-The-Loop)
// ---------------------------------------------------------------------------

// HitlInterruptEvent is emitted when the agent pauses awaiting human approval.
// Submit decisions via ChatService.Decide to resume the run.
type HitlInterruptEvent struct {
	ID          string                 `json:"id"`
	ActionName  string                 `json:"action_name"`
	ActionArgs  map[string]interface{} `json:"action_args"`
	Description string                 `json:"description"`
}

func (*HitlInterruptEvent) Channel() string { return "hitl" }
func (*HitlInterruptEvent) Type() string    { return "hitl_interrupt" }

// ---------------------------------------------------------------------------
// Error / cancellation
// ---------------------------------------------------------------------------

// ErrorEvent is emitted when the run terminates due to an unrecoverable error.
type ErrorEvent struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

func (*ErrorEvent) Channel() string { return "error" }
func (*ErrorEvent) Type() string    { return "error" }

// CancelledEvent is emitted after a stop request completes.
type CancelledEvent struct {
	Reason string `json:"reason"`
}

func (*CancelledEvent) Channel() string { return "error" }
func (*CancelledEvent) Type() string    { return "cancelled" }

// ---------------------------------------------------------------------------
// Heartbeat
// ---------------------------------------------------------------------------

// HeartbeatEvent is a periodic keep-alive; applications may safely ignore it.
type HeartbeatEvent struct{}

func (*HeartbeatEvent) Channel() string { return "heartbeat" }
func (*HeartbeatEvent) Type() string    { return "" }
