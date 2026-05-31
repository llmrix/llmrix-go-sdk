package llmrix

// This file re-exports the most commonly used types from the model/ and streaming/
// sub-packages. Callers can use e.g. llmrix.Conversation directly without an
// additional import, while packages that prefer explicit imports can still use
// model.Conversation / streaming.MessageChunkEvent.

import (
	"github.com/llmrix/llmrix-go-sdk/model"
	"github.com/llmrix/llmrix-go-sdk/streaming"
)

// ---------------------------------------------------------------------------
// Model re-exports
// ---------------------------------------------------------------------------

type (
	Conversation              = model.Conversation
	ConversationCreateRequest = model.ConversationCreateRequest
	ConversationUpdateRequest = model.ConversationUpdateRequest
	Message                   = model.Message
	ToolCall                  = model.ToolCall

	ChatRequest       = model.ChatRequest
	HitlDecision      = model.HitlDecision
	HitlDecideRequest = model.HitlDecideRequest

	CronTask          = model.CronTask
	CronCreateRequest = model.CronCreateRequest
	CronUpdateRequest = model.CronUpdateRequest

	Agent              = model.Agent
	AgentCreateRequest = model.AgentCreateRequest
	AgentUpdateRequest = model.AgentUpdateRequest
	Mate               = model.Mate
	SaveMatesRequest   = model.SaveMatesRequest
)

// HITL decision factory shortcuts.
var (
	Approve = model.Approve
	Reject  = model.Reject
	Modify  = model.Modify
)

// ---------------------------------------------------------------------------
// Streaming re-exports
// ---------------------------------------------------------------------------

type (
	StreamEvent = streaming.StreamEvent
	EventHandler = streaming.EventHandler

	RunStartEvent      = streaming.RunStartEvent
	RunEndEvent        = streaming.RunEndEvent
	MessageChunkEvent  = streaming.MessageChunkEvent
	ToolStartEvent     = streaming.ToolStartEvent
	ToolEndEvent       = streaming.ToolEndEvent
	SubagentStartEvent = streaming.SubagentStartEvent
	SubagentEndEvent   = streaming.SubagentEndEvent
	HitlInterruptEvent = streaming.HitlInterruptEvent
	ErrorEvent         = streaming.ErrorEvent
	CancelledEvent     = streaming.CancelledEvent
	HeartbeatEvent     = streaming.HeartbeatEvent
)

// RegisterEvent adds a custom channel:type → factory mapping to the SSE registry.
// Allows handling server-extension events without forking the SDK.
var RegisterEvent = streaming.RegisterEvent
