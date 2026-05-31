package llmrix

import (
	"context"
	"net/http"

	"github.com/llmrix/llmrix-go-sdk/internal/transport"
	"github.com/llmrix/llmrix-go-sdk/model"
)

// CronService provides operations on the Cron Tasks API.
// Obtain via LlmrixClient.Cron().
type CronService struct {
	t *transport.Transport
}

// List returns all cron tasks.
func (s *CronService) List(ctx context.Context) ([]model.CronTask, error) {
	var out []model.CronTask
	err := s.t.DoUnwrap(ctx, http.MethodGet, transport.PathCronTasks, nil, "tasks", &out)
	return out, err
}

// Create creates a new cron task.
func (s *CronService) Create(ctx context.Context, req model.CronCreateRequest) (*model.CronTask, error) {
	var out model.CronTask
	err := s.t.Do(ctx, http.MethodPost, transport.PathCronTasks, req, &out)
	return &out, err
}

// Get retrieves a single cron task by its UUID.
func (s *CronService) Get(ctx context.Context, taskID string) (*model.CronTask, error) {
	var out model.CronTask
	err := s.t.Do(ctx, http.MethodGet, transport.PathCronTask(taskID), nil, &out)
	return &out, err
}

// Update updates a cron task's properties.
func (s *CronService) Update(ctx context.Context, taskID string, req model.CronUpdateRequest) (*model.CronTask, error) {
	var out model.CronTask
	err := s.t.Do(ctx, http.MethodPut, transport.PathCronTask(taskID), req, &out)
	return &out, err
}

// Pause pauses an active cron task.
func (s *CronService) Pause(ctx context.Context, taskID string) (*model.CronTask, error) {
	var out model.CronTask
	err := s.t.Do(ctx, http.MethodPost, transport.PathCronTaskPause(taskID), struct{}{}, &out)
	return &out, err
}

// Resume resumes a paused cron task.
func (s *CronService) Resume(ctx context.Context, taskID string) (*model.CronTask, error) {
	var out model.CronTask
	err := s.t.Do(ctx, http.MethodPost, transport.PathCronTaskResume(taskID), struct{}{}, &out)
	return &out, err
}

// Delete permanently deletes a cron task.
func (s *CronService) Delete(ctx context.Context, taskID string) error {
	return s.t.Do(ctx, http.MethodDelete, transport.PathCronTask(taskID), nil, nil)
}
