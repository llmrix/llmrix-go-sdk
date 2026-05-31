package llmrix

import (
	"context"
	"net/http"

	"github.com/llmrix/llmrix-go-sdk/internal/transport"
	"github.com/llmrix/llmrix-go-sdk/model"
)

// AgentService provides operations on the Agents API.
// Obtain via LlmrixClient.Agents().
// All methods return HTTP 503 when the server runs in native/standalone mode.
type AgentService struct {
	t *transport.Transport
}

// List returns all agents.
func (s *AgentService) List(ctx context.Context) ([]model.Agent, error) {
	var out []model.Agent
	err := s.t.DoUnwrap(ctx, http.MethodGet, transport.PathAgents, nil, "agents", &out)
	return out, err
}

// Create creates a new agent.
func (s *AgentService) Create(ctx context.Context, req model.AgentCreateRequest) (*model.Agent, error) {
	var out model.Agent
	err := s.t.DoUnwrap(ctx, http.MethodPost, transport.PathAgents, req, "agent", &out)
	return &out, err
}

// Get retrieves a single agent by numeric ID.
func (s *AgentService) Get(ctx context.Context, agentID int) (*model.Agent, error) {
	var out model.Agent
	err := s.t.DoUnwrap(ctx, http.MethodGet, transport.PathAgent(agentID), nil, "agent", &out)
	return &out, err
}

// Update updates an agent's properties. Only non-nil pointer fields are sent.
func (s *AgentService) Update(ctx context.Context, agentID int, req model.AgentUpdateRequest) (*model.Agent, error) {
	var out model.Agent
	err := s.t.DoUnwrap(ctx, http.MethodPatch, transport.PathAgent(agentID), req, "agent", &out)
	return &out, err
}

// Delete permanently deletes an agent.
func (s *AgentService) Delete(ctx context.Context, agentID int) error {
	return s.t.Do(ctx, http.MethodDelete, transport.PathAgent(agentID), nil, nil)
}

// ListMates returns all mates assigned to a team agent.
func (s *AgentService) ListMates(ctx context.Context, agentID int) ([]model.Mate, error) {
	var out []model.Mate
	err := s.t.DoUnwrap(ctx, http.MethodGet, transport.PathAgentMates(agentID), nil, "mates", &out)
	return out, err
}

// SaveMates atomically replaces all mates of a team agent.
func (s *AgentService) SaveMates(ctx context.Context, agentID int, mates []model.Mate) ([]model.Mate, error) {
	var out []model.Mate
	err := s.t.DoUnwrap(ctx, http.MethodPost, transport.PathAgentMates(agentID),
		model.SaveMatesRequest{Mates: mates}, "mates", &out)
	return out, err
}
