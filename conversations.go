package llmrix

import (
	"context"
	"fmt"
	"net/http"

	"github.com/llmrix/llmrix-go-sdk/internal/transport"
	"github.com/llmrix/llmrix-go-sdk/model"
)

// ConversationsService provides operations on the Conversations API.
// Obtain via LlmrixClient.Conversations().
type ConversationsService struct {
	t *transport.Transport
}

// Create creates a new conversation.
func (s *ConversationsService) Create(ctx context.Context, req model.ConversationCreateRequest) (*model.Conversation, error) {
	var out model.Conversation
	err := s.t.Do(ctx, http.MethodPost, transport.PathConversations, req, &out)
	return &out, err
}

// List returns a page of conversations ordered by last activity.
// Pass lastID=0 to start from the beginning; use the Seq of the last item to advance.
func (s *ConversationsService) List(ctx context.Context, lastID int64, size int) (*PageResult[model.Conversation], error) {
	path := fmt.Sprintf("%s?lastId=%d&size=%d", transport.PathConversations, lastID, size)
	var out PageResult[model.Conversation]
	err := s.t.Do(ctx, http.MethodGet, path, nil, &out)
	return &out, err
}

// Get retrieves a single conversation by ID.
func (s *ConversationsService) Get(ctx context.Context, id string) (*model.Conversation, error) {
	var out model.Conversation
	err := s.t.Do(ctx, http.MethodGet, transport.PathConversation(id), nil, &out)
	return &out, err
}

// Update updates a conversation's metadata (title).
func (s *ConversationsService) Update(ctx context.Context, id string, req model.ConversationUpdateRequest) (*model.Conversation, error) {
	var out model.Conversation
	err := s.t.Do(ctx, http.MethodPatch, transport.PathConversation(id), req, &out)
	return &out, err
}

// Delete permanently deletes a conversation and all of its messages.
func (s *ConversationsService) Delete(ctx context.Context, id string) error {
	return s.t.Do(ctx, http.MethodDelete, transport.PathConversation(id), nil, nil)
}

// Messages returns a page of messages in the given conversation.
// Pass lastID=0 to start from the beginning.
func (s *ConversationsService) Messages(ctx context.Context, id string, lastID int64, size int) (*PageResult[model.Message], error) {
	path := fmt.Sprintf("%s?lastId=%d&size=%d", transport.PathMessages(id), lastID, size)
	var out PageResult[model.Message]
	err := s.t.Do(ctx, http.MethodGet, path, nil, &out)
	return &out, err
}
