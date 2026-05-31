package llmrix

import (
	"context"
	"net/http"

	"github.com/llmrix/llmrix-go-sdk/internal/apierr"
	"github.com/llmrix/llmrix-go-sdk/internal/transport"
	"github.com/llmrix/llmrix-go-sdk/model"
	"github.com/llmrix/llmrix-go-sdk/streaming"
)

// ChatService provides streaming chat, stop, and HITL-decision operations
// scoped to a single conversation. Obtain via LlmrixClient.Chat(convID).
type ChatService struct {
	t      *transport.Transport
	convID string
}

// Send sends a plain-text message and streams the agent's response.
// handler is called synchronously for each received StreamEvent; return a
// non-nil error from handler to abort streaming early.
//
//	err := client.Chat(convID).Send(ctx, "Hello!", func(e llmrix.StreamEvent) error {
//	    if chunk, ok := e.(*llmrix.MessageChunkEvent); ok {
//	        fmt.Print(chunk.Content)
//	    }
//	    return nil
//	})
func (s *ChatService) Send(ctx context.Context, message string, handler streaming.EventHandler) error {
	return s.SendRequest(ctx, model.ChatRequest{Message: message}, handler)
}

// SendRequest sends a fully-specified ChatRequest and streams the response.
// Use this overload to supply AgentID, Metadata, or HITL decisions inline.
func (s *ChatService) SendRequest(ctx context.Context, req model.ChatRequest, handler streaming.EventHandler) error {
	return s.t.Stream(ctx, transport.PathChat(s.convID), req, handler)
}

// Stop requests the server to cancel the currently running chat turn.
// The in-flight Send call will receive a CancelledEvent before the stream closes.
func (s *ChatService) Stop(ctx context.Context) error {
	return s.t.Do(ctx, http.MethodPost, transport.PathChatStop(s.convID), struct{}{}, nil)
}

// Decide submits HITL decisions to resume a paused agent run.
// Call this after receiving a HitlInterruptEvent.
func (s *ChatService) Decide(ctx context.Context, decisions []model.HitlDecision) error {
	if len(decisions) == 0 {
		return &apierr.LlmrixError{Message: "decisions must not be empty"}
	}
	return s.t.Do(ctx, http.MethodPost, transport.PathChatHitlDecide(s.convID),
		model.HitlDecideRequest{Decisions: decisions}, nil)
}
