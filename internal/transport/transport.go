// Package transport provides the low-level HTTP and SSE transport used by all services.
// It is internal to the SDK and must not be imported by external consumers.
package transport

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/llmrix/llmrix-go-sdk/internal/apierr"
	"github.com/llmrix/llmrix-go-sdk/streaming"
)

// Transport handles all low-level HTTP operations: building requests, executing
// them with the correct auth header, and mapping error responses to typed errors.
// Regular requests use a timeout-bounded client; SSE streams use a no-timeout client.
type Transport struct {
	baseURL   string
	apiKey    string
	client    *http.Client // regular requests (has timeout)
	sseClient *http.Client // SSE streaming  (no timeout — context controls lifetime)
}

// New creates a Transport. sseClient should share Transport with regularClient so
// both use the same connection pool.
func New(baseURL, apiKey string, regularClient, sseClient *http.Client) *Transport {
	return &Transport{
		baseURL:   baseURL,
		apiKey:    apiKey,
		client:    regularClient,
		sseClient: sseClient,
	}
}

// Do executes a request and deserialises a successful 2xx JSON body into dest.
// Pass dest=nil when no response body is needed.
func (t *Transport) Do(ctx context.Context, method, path string, body, dest interface{}) error {
	req, err := t.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return apierr.New("request failed", err)
	}
	defer resp.Body.Close()

	if err := t.checkStatus(resp); err != nil {
		return err
	}
	if dest == nil {
		return nil
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return apierr.New("failed to read response body", err)
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return apierr.New(fmt.Sprintf("failed to deserialize response: %.200s", string(raw)), err)
	}
	return nil
}

// DoUnwrap executes a request, parses the JSON envelope, and extracts the value
// at key into dest. Used for responses shaped like {"agent": {...}}.
func (t *Transport) DoUnwrap(ctx context.Context, method, path string, body interface{}, key string, dest interface{}) error {
	req, err := t.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return apierr.New("request failed", err)
	}
	defer resp.Body.Close()

	if err := t.checkStatus(resp); err != nil {
		return err
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return apierr.New("failed to read response body", err)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return apierr.New("failed to parse response envelope", err)
	}
	sub, ok := envelope[key]
	if !ok {
		// Key absent — if dest is a slice pointer, return empty rather than error.
		rv := reflect.ValueOf(dest)
		if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Slice {
			return nil
		}
		return apierr.New(fmt.Sprintf("response missing key %q", key), nil)
	}
	if err := json.Unmarshal(sub, dest); err != nil {
		return apierr.New(fmt.Sprintf("failed to deserialize field %q", key), err)
	}
	return nil
}

// Stream POSTs to path using the no-timeout SSE client and dispatches each
// parsed event to handler. Blocks until the server closes the stream.
func (t *Transport) Stream(ctx context.Context, path string, body interface{}, handler streaming.EventHandler) error {
	req, err := t.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := t.sseClient.Do(req)
	if err != nil {
		return apierr.New("SSE request failed", err)
	}
	defer resp.Body.Close()

	if err := t.checkStatus(resp); err != nil {
		return err
	}
	return parseSse(resp.Body, handler)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func (t *Transport) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, apierr.New("failed to serialize request body", err)
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, t.baseURL+path, bodyReader)
	if err != nil {
		return nil, apierr.New("failed to build request", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if t.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+t.apiKey)
	}
	return req, nil
}

func (t *Transport) checkStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	msg := extractErrorMessage(body, resp.StatusCode)
	return apierr.NewAPI(resp.StatusCode, msg, string(body))
}

func extractErrorMessage(body []byte, code int) string {
	if len(body) == 0 {
		return fmt.Sprintf("server returned HTTP %d", code)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err == nil {
		for _, key := range []string{"message", "detail", "error", "msg"} {
			if v, ok := m[key]; ok {
				if s, ok := v.(string); ok {
					return s
				}
			}
		}
	}
	s := string(body)
	if len(s) > 200 {
		return s[:200] + "..."
	}
	return s
}

// ---------------------------------------------------------------------------
// SSE parser  (moved here from root; uses streaming.EventHandler)
// ---------------------------------------------------------------------------

func parseSse(r io.Reader, handler streaming.EventHandler) error {
	scanner := bufio.NewScanner(r)
	var (
		eventName string
		dataBuf   strings.Builder
	)

	dispatch := func() error {
		data := dataBuf.String()
		dataBuf.Reset()
		if data == "" {
			return nil
		}
		return streaming.DispatchFrame(eventName, data, handler)
	}

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "":
			if err := dispatch(); err != nil {
				return err
			}
			eventName = ""
		case strings.HasPrefix(line, "event:"):
			eventName = strings.TrimSpace(line[len("event:"):])
		case strings.HasPrefix(line, "data:"):
			frag := line[len("data:"):]
			if len(frag) > 0 && frag[0] == ' ' {
				frag = frag[1:]
			}
			if dataBuf.Len() > 0 {
				dataBuf.WriteByte('\n')
			}
			dataBuf.WriteString(frag)
		}
	}
	if err := dispatch(); err != nil {
		return err
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		return apierr.New("SSE stream read error", err)
	}
	return nil
}
