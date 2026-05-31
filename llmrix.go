// Package llmrix is the official Go SDK for the llmrix AI Agent Platform.
//
// Create a client and start using it:
//
//	client := llmrix.NewClient(
//	    llmrix.WithBaseURL("http://localhost:8899"),
//	    llmrix.WithAPIKey("sk-xxx"),
//	)
//
//	conv, err := client.Conversations().Create(ctx,
//	    model.ConversationCreateRequest{Title: "My chat"})
//
//	err = client.Chat(conv.ID).Send(ctx, "Hello!", func(e llmrix.StreamEvent) error {
//	    if chunk, ok := e.(*llmrix.MessageChunkEvent); ok {
//	        fmt.Print(chunk.Content)
//	    }
//	    return nil
//	})
package llmrix

import (
	"net/http"
	"strings"
	"time"

	"github.com/llmrix/llmrix-go-sdk/internal/transport"
)

// LlmrixClient is the main entry point for the Llmrix Go SDK.
// It is safe for concurrent use; the underlying HTTP client uses a shared connection pool.
type LlmrixClient struct {
	t             *transport.Transport
	conversations *ConversationsService
	cron          *CronService
	agents        *AgentService
}

// NewClient creates a new LlmrixClient with the given options.
// WithBaseURL is required; all other options have sensible defaults.
func NewClient(opts ...Option) *LlmrixClient {
	cfg := &clientConfig{timeout: 60 * time.Second}
	for _, o := range opts {
		o(cfg)
	}
	if cfg.baseURL == "" {
		panic("llmrix: WithBaseURL is required")
	}
	cfg.baseURL = strings.TrimRight(cfg.baseURL, "/")

	regular := cfg.httpClient
	if regular == nil {
		regular = &http.Client{Timeout: cfg.timeout}
	}
	// SSE client shares the connection pool but has no overall timeout;
	// the caller's context controls stream lifetime.
	sse := &http.Client{Transport: regular.Transport, Timeout: 0}

	t := transport.New(cfg.baseURL, cfg.apiKey, regular, sse)
	c := &LlmrixClient{t: t}
	c.conversations = &ConversationsService{t: t}
	c.cron = &CronService{t: t}
	c.agents = &AgentService{t: t}
	return c
}

// Conversations returns the service for managing conversations and message history.
func (c *LlmrixClient) Conversations() *ConversationsService { return c.conversations }

// Cron returns the service for managing scheduled cron tasks.
func (c *LlmrixClient) Cron() *CronService { return c.cron }

// Agents returns the service for managing agents and their mates.
// Note: available in cloud mode only; returns HTTP 503 in native/standalone mode.
func (c *LlmrixClient) Agents() *AgentService { return c.agents }

// Chat returns a ChatService scoped to the given conversation ID.
func (c *LlmrixClient) Chat(convID string) *ChatService {
	return &ChatService{t: c.t, convID: convID}
}

// ---------------------------------------------------------------------------
// Pagination helper (generic — defined here because Go <1.24 cannot alias generics)
// ---------------------------------------------------------------------------

// PageResult is the generic paged list returned by all list endpoints.
type PageResult[T any] struct {
	Items   []T  `json:"items"`
	HasMore bool `json:"has_more"`
}

// ---------------------------------------------------------------------------
// Functional options
// ---------------------------------------------------------------------------

type clientConfig struct {
	baseURL    string
	apiKey     string
	timeout    time.Duration
	httpClient *http.Client
}

// Option is a functional option for NewClient.
type Option func(*clientConfig)

// WithBaseURL sets the Llmrix server base URL (e.g. "http://localhost:8899").
// Trailing slashes are stripped automatically. Required.
func WithBaseURL(url string) Option { return func(c *clientConfig) { c.baseURL = url } }

// WithAPIKey sets the API key sent in the Authorization: Bearer header.
// Omit if the server has no auth configured.
func WithAPIKey(key string) Option { return func(c *clientConfig) { c.apiKey = key } }

// WithTimeout sets the HTTP request timeout for non-streaming requests (default: 60s).
// SSE streaming calls respect context deadlines instead.
func WithTimeout(d time.Duration) Option { return func(c *clientConfig) { c.timeout = d } }

// WithHTTPClient replaces the default *http.Client for non-streaming requests.
// The SSE client always has Timeout=0 but shares the same Transport.
func WithHTTPClient(hc *http.Client) Option { return func(c *clientConfig) { c.httpClient = hc } }
