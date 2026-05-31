package streaming

import (
	"encoding/json"
	"fmt"
)

// EventHandler is a callback invoked once per parsed StreamEvent.
// Return a non-nil error to abort streaming early.
type EventHandler func(StreamEvent) error

// ---------------------------------------------------------------------------
// Registry  (extensible, parallel to Java SseStreamParser.REGISTRY)
// ---------------------------------------------------------------------------

// registry maps "channel:type" → factory that returns a zero-value event.
var registry = map[string]func() StreamEvent{
	"lifecycle:run_start":    func() StreamEvent { return &RunStartEvent{} },
	"lifecycle:run_end":      func() StreamEvent { return &RunEndEvent{} },
	"messages:message_chunk": func() StreamEvent { return &MessageChunkEvent{} },
	"tools:tool_start":       func() StreamEvent { return &ToolStartEvent{} },
	"tools:tool_end":         func() StreamEvent { return &ToolEndEvent{} },
	"tools:subagent_start":   func() StreamEvent { return &SubagentStartEvent{} },
	"tools:subagent_end":     func() StreamEvent { return &SubagentEndEvent{} },
	"hitl:hitl_interrupt":    func() StreamEvent { return &HitlInterruptEvent{} },
	"error:error":            func() StreamEvent { return &ErrorEvent{} },
	"error:cancelled":        func() StreamEvent { return &CancelledEvent{} },
	"heartbeat:":             func() StreamEvent { return &HeartbeatEvent{} },
}

// RegisterEvent adds a custom channel:type → factory mapping, enabling callers
// to handle server-extension events without forking the SDK.
//
//	streaming.RegisterEvent("audit", "log_entry", func() streaming.StreamEvent {
//	    return &MyAuditEvent{}
//	})
func RegisterEvent(channel, typ string, factory func() StreamEvent) {
	registry[channel+":"+typ] = factory
}

// ---------------------------------------------------------------------------
// Frame dispatcher  (called from internal/transport after raw SSE parsing)
// ---------------------------------------------------------------------------

// DispatchFrame parses a single SSE data frame and calls handler with the
// resolved typed event. Unknown channels/types are silently skipped.
// Exported so internal/transport can call it without duplicating logic.
func DispatchFrame(eventName, data string, handler EventHandler) error {
	// Bare ping / heartbeat
	if data == ":ping" || eventName == "heartbeat" {
		return handler(&HeartbeatEvent{})
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil // malformed frame — skip
	}

	channel := jsonString(raw["channel"])
	if channel == "" {
		channel = eventName
	}
	typ := jsonString(raw["type"])

	key := channel + ":" + typ
	factory, ok := registry[key]
	if !ok {
		// Channel-only fallback (e.g. "heartbeat:")
		factory, ok = registry[channel+":"]
		if !ok {
			return nil // unknown event — skip
		}
	}

	evt := factory()
	if err := json.Unmarshal([]byte(data), evt); err != nil {
		return fmt.Errorf("llmrix: failed to deserialize %s event: %w", key, err)
	}
	return handler(evt)
}

func jsonString(raw json.RawMessage) string {
	if raw == nil {
		return ""
	}
	var s string
	_ = json.Unmarshal(raw, &s)
	return s
}
