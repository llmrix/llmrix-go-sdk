# llmrix-go-sdk

Official Go SDK for [llmrix](https://github.com/llmrix/llmrix-go-sdk) — AI Agent Platform.

- Go 1.21+
- Zero external dependencies (stdlib `net/http` only)
- Context-aware — every call accepts a `context.Context`
- SSE streaming with a clean callback interface
- Thread-safe client; lightweight per-chat service objects

---

## Installation

```bash
go get github.com/llmrix/llmrix-go-sdk
```

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/llmrix/llmrix-go-sdk"
    "github.com/llmrix/llmrix-go-sdk/model"
)

func main() {
    ctx := context.Background()

    // 1. Create the client (thread-safe, reuse across goroutines)
    client := llmrix.NewClient(
        llmrix.WithBaseURL("https://www.llmrix.com"),
        llmrix.WithAPIKey("sk-xxx"),
    )

    // 2. Create a conversation
    conv, err := client.Conversations().Create(ctx,
        model.ConversationCreateRequest{Title: "My Chat"})
    if err != nil {
        panic(err)
    }

    // 3. Stream a chat turn
    err = client.Chat(conv.ID).Send(ctx, "Hello, llmrix!", func(e llmrix.StreamEvent) error {
        if chunk, ok := e.(*llmrix.MessageChunkEvent); ok {
            fmt.Print(chunk.Content)
        }
        return nil
    })
    if err != nil {
        panic(err)
    }
}
```

---

## API Reference

### `llmrix.NewClient`

```go
func NewClient(opts ...Option) *LlmrixClient
```

| Option | Description |
|--------|-------------|
| `WithBaseURL(url string)` | Server base URL — **required** |
| `WithAPIKey(key string)` | Bearer token for `Authorization` header |
| `WithTimeout(d time.Duration)` | Non-streaming request timeout (default: 60 s) |
| `WithHTTPClient(hc *http.Client)` | Replace the default HTTP client |

---

### `ConversationsService`

Obtained via `client.Conversations()`.

| Method | Returns | Description |
|--------|---------|-------------|
| `Create(ctx, req)` | `(*Conversation, error)` | Create a new conversation |
| `List(ctx, lastID, size)` | `(*PageResult[Conversation], error)` | Cursor-paginated list |
| `Get(ctx, id)` | `(*Conversation, error)` | Fetch a single conversation |
| `Update(ctx, id, req)` | `(*Conversation, error)` | Update title / agent |
| `Delete(ctx, id)` | `error` | Permanently delete |
| `Messages(ctx, id, lastID, size)` | `(*PageResult[Message], error)` | Paginated message history |

Cursor pagination example:

```go
page, _ := client.Conversations().List(ctx, 0, 20)
for page.HasMore {
    last := page.Items[len(page.Items)-1]
    page, _ = client.Conversations().List(ctx, last.Seq, 20)
}
```

---

### `ChatService`

Obtained via `client.Chat(convID)`.

| Method | Returns | Description |
|--------|---------|-------------|
| `Send(ctx, message, handler)` | `error` | Send plain-text message, stream response |
| `SendRequest(ctx, req, handler)` | `error` | Send full `ChatRequest`, stream response |
| `Stop(ctx)` | `error` | Cancel the current run |
| `Decide(ctx, decisions)` | `error` | Submit HITL decisions |

#### Stream handler

```go
type StreamHandler func(StreamEvent) error
```

If the handler returns a non-nil error the stream is aborted.

---

### Stream Event Types

All events implement the `StreamEvent` interface. Use type assertions or a type switch:

```go
switch e := event.(type) {
case *llmrix.MessageChunkEvent:
    fmt.Print(e.Content)
case *llmrix.ToolStartEvent:
    fmt.Printf("Tool: %s\n", e.Name)
case *llmrix.HitlInterruptEvent:
    _ = client.Chat(conv.ID).Decide(ctx, []model.HitlDecision{model.HitlApprove()})
case *llmrix.RunEndEvent:
    fmt.Println("\n[done]")
case *llmrix.ErrorEvent:
    fmt.Fprintf(os.Stderr, "error: %s\n", e.Message)
}
```

| Type | Channel | Description |
|------|---------|-------------|
| `RunStartEvent` | `lifecycle` | Agent run started |
| `RunEndEvent` | `lifecycle` | Agent run completed |
| `MessageChunkEvent` | `messages` | Incremental text chunk |
| `ToolStartEvent` | `tools` | Tool invocation started |
| `ToolEndEvent` | `tools` | Tool invocation finished |
| `SubagentStartEvent` | `tools` | Sub-agent spawned |
| `SubagentEndEvent` | `tools` | Sub-agent finished |
| `HitlInterruptEvent` | `hitl` | Human approval required |
| `ErrorEvent` | `error` | Run failed with error |
| `CancelledEvent` | `error` | Run was cancelled |
| `HeartbeatEvent` | `heartbeat` | Keep-alive (silently dropped) |

---

### Error Handling

All service methods return a standard Go `error`. SDK-specific errors can be inspected
with `errors.As`:

```go
import "github.com/llmrix/llmrix-go-sdk"

var apiErr *llmrix.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.Body)
}
```

| Type | Description |
|------|-------------|
| `*llmrix.APIError` | Non-2xx HTTP response; exposes `StatusCode` and `Body` |
| `*llmrix.AuthError` | HTTP 401 / 403 |

---

## Advanced Usage

### HITL (Human-in-the-Loop)

```go
chat := client.Chat(conv.ID)
_ = chat.Send(ctx, "Delete everything in /tmp", func(e llmrix.StreamEvent) error {
    if interrupt, ok := e.(*llmrix.HitlInterruptEvent); ok {
        fmt.Printf("Approve action: %s? (y/n): ", interrupt.ActionName)
        var ans string
        fmt.Scan(&ans)
        decision := model.HitlReject("user declined")
        if ans == "y" {
            decision = model.HitlApprove()
        }
        return chat.Decide(ctx, []model.HitlDecision{decision})
    }
    return nil
})
```

### Collecting the full response

```go
func runToCompletion(client *llmrix.LlmrixClient, convID, message string) (string, error) {
    var buf strings.Builder
    err := client.Chat(convID).Send(context.Background(), message, func(e llmrix.StreamEvent) error {
        if chunk, ok := e.(*llmrix.MessageChunkEvent); ok {
            buf.WriteString(chunk.Content)
        }
        return nil
    })
    return buf.String(), err
}
```

---

## Building from Source

```bash
git clone https://github.com/llmrix/llmrix-go-sdk.git
cd llmrix-go-sdk
go build ./...
go test ./...
```

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

---

## License

Apache License 2.0 — see the [LICENSE](LICENSE) file.
