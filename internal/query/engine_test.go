package query

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/user-name/cc-cli-go/internal/api"
	"github.com/user-name/cc-cli-go/internal/permission"
	"github.com/user-name/cc-cli-go/internal/tools"
	"github.com/user-name/cc-cli-go/internal/types"
)

type scriptStreamer struct {
	mu    sync.Mutex
	turns [][]api.StreamEvent
	next  int
	reqs  []*api.Request
}

func (s *scriptStreamer) Stream(_ context.Context, req *api.Request) (<-chan api.StreamEvent, error) {
	raw, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var snap api.Request
	if err := json.Unmarshal(raw, &snap); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.reqs = append(s.reqs, &snap)
	var events []api.StreamEvent
	if s.next < len(s.turns) {
		events = s.turns[s.next]
		s.next++
	}
	s.mu.Unlock()

	ch := make(chan api.StreamEvent, len(events))
	for _, event := range events {
		ch <- event
	}
	close(ch)
	return ch, nil
}

func (s *scriptStreamer) bodies() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, len(s.reqs))
	for i, req := range s.reqs {
		raw, _ := json.Marshal(req)
		out[i] = string(raw)
	}
	return out
}

type execTool struct {
	name  string
	out   string
	mu    sync.Mutex
	calls int
}

func (t *execTool) Name() string        { return t.name }
func (t *execTool) Description() string { return "test tool" }
func (t *execTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{"type": "object"}
}
func (t *execTool) IsEnabled() bool                               { return true }
func (t *execTool) IsReadOnly(map[string]interface{}) bool        { return true }
func (t *execTool) IsConcurrencySafe(map[string]interface{}) bool { return true }
func (t *execTool) UserFacingName(map[string]interface{}) string  { return t.name }
func (t *execTool) Execute(context.Context, map[string]interface{}, *tools.ToolContext) (*tools.ToolResult, error) {
	t.mu.Lock()
	t.calls++
	t.mu.Unlock()
	return &tools.ToolResult{Content: t.out}, nil
}

func (t *execTool) Calls() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.calls
}

func toolUseTurn(id, name string, input map[string]interface{}) []api.StreamEvent {
	msg, _ := json.Marshal(map[string]string{"id": "msg_" + id, "role": "assistant"})
	block, _ := json.Marshal(types.ContentBlock{Type: "tool_use", ID: id, Name: name, Input: input})
	delta, _ := json.Marshal(map[string]interface{}{
		"type":  "message_delta",
		"delta": map[string]string{"stop_reason": "tool_use"},
	})
	return []api.StreamEvent{
		{Type: "message_start", Message: msg},
		{Type: "content_block_start", ContentBlock: block},
		{Type: "content_block_stop"},
		{Type: "message_delta", Delta: delta},
		{Type: "message_stop"},
	}
}

func endTurn() []api.StreamEvent {
	msg, _ := json.Marshal(map[string]string{"id": "msg_end", "role": "assistant"})
	block, _ := json.Marshal(types.ContentBlock{Type: "text", Text: "done"})
	delta, _ := json.Marshal(map[string]interface{}{
		"type":  "message_delta",
		"delta": map[string]string{"stop_reason": "end_turn"},
	})
	return []api.StreamEvent{
		{Type: "message_start", Message: msg},
		{Type: "content_block_start", ContentBlock: block},
		{Type: "content_block_stop"},
		{Type: "message_delta", Delta: delta},
		{Type: "message_stop"},
	}
}

func drain(events <-chan StreamEvent) {
	for range events {
	}
}

func runQuery(t *testing.T, eng *Engine, params QueryParams) QueryResult {
	t.Helper()
	events, results := eng.Query(context.Background(), params)
	done := make(chan struct{})
	go func() {
		drain(events)
		close(done)
	}()
	res := <-results
	<-done
	return res
}

func TestQuery_FollowUpRequestContainsToolResult(t *testing.T) {
	tool := &execTool{name: "Echo", out: "tool-output-42"}
	reg := tools.NewRegistry()
	reg.Register(tool)
	script := &scriptStreamer{turns: [][]api.StreamEvent{
		toolUseTurn("toolu_1", "Echo", map[string]interface{}{"q": "hi"}),
		endTurn(),
	}}
	checker := permission.NewChecker(permission.ModeAccept)
	eng := &Engine{client: script, toolReg: reg}

	res := runQuery(t, eng, QueryParams{
		Messages:          []*types.Message{types.NewUserMessage("hi")},
		Tools:             []tools.Tool{tool},
		Model:             "test-model",
		MaxTokens:         32,
		MaxTurns:          4,
		PermissionChecker: checker,
	})

	if res.Error != nil {
		t.Fatalf("query error: %v", res.Error)
	}
	if res.Reason != "completed" {
		t.Fatalf("reason = %q", res.Reason)
	}
	if tool.Calls() != 1 {
		t.Fatalf("Execute calls = %d", tool.Calls())
	}
	bodies := script.bodies()
	if len(bodies) < 2 {
		t.Fatalf("model requests = %d", len(bodies))
	}
	if !strings.Contains(bodies[1], "tool-output-42") {
		t.Fatalf("follow-up request missing tool output: %s", bodies[1])
	}
	if !strings.Contains(bodies[1], "toolu_1") {
		t.Fatalf("follow-up request missing tool use id: %s", bodies[1])
	}
}

func TestQuery_DenyAndAskDoNotExecute(t *testing.T) {
	cases := []struct {
		name    string
		rules   []permission.Rule
		snippet string
	}{
		{
			name:    "deny",
			rules:   []permission.Rule{{ToolName: "Echo", Pattern: "*", Behavior: permission.BehaviorDeny}},
			snippet: "Permission denied",
		},
		{
			name:    "ask",
			rules:   nil,
			snippet: "Permission required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tool := &execTool{name: "Echo", out: "should-not-run"}
			reg := tools.NewRegistry()
			reg.Register(tool)
			script := &scriptStreamer{turns: [][]api.StreamEvent{
				toolUseTurn("toolu_perm", "Echo", map[string]interface{}{"q": "x"}),
				endTurn(),
			}}
			checker := permission.NewChecker(permission.ModeDefault)
			if tc.rules != nil {
				checker.SetRules(tc.rules)
			}
			eng := &Engine{client: script, toolReg: reg}
			res := runQuery(t, eng, QueryParams{
				Messages:          []*types.Message{types.NewUserMessage("hi")},
				Tools:             []tools.Tool{tool},
				Model:             "test-model",
				MaxTokens:         32,
				MaxTurns:          4,
				PermissionChecker: checker,
			})
			if res.Error != nil {
				t.Fatalf("query error: %v", res.Error)
			}
			if tool.Calls() != 0 {
				t.Fatalf("Execute calls = %d", tool.Calls())
			}
			bodies := script.bodies()
			if len(bodies) < 2 {
				t.Fatalf("model requests = %d", len(bodies))
			}
			if strings.Contains(bodies[1], "should-not-run") {
				t.Fatalf("tool output leaked into follow-up: %s", bodies[1])
			}
			if !strings.Contains(bodies[1], tc.snippet) {
				t.Fatalf("follow-up missing %q: %s", tc.snippet, bodies[1])
			}
		})
	}
}

func TestQuery_MaxTurnsStopsToolLoop(t *testing.T) {
	tool := &execTool{name: "Echo", out: "again"}
	reg := tools.NewRegistry()
	reg.Register(tool)
	script := &scriptStreamer{turns: [][]api.StreamEvent{
		toolUseTurn("toolu_a", "Echo", map[string]interface{}{"q": "1"}),
		toolUseTurn("toolu_b", "Echo", map[string]interface{}{"q": "2"}),
		toolUseTurn("toolu_c", "Echo", map[string]interface{}{"q": "3"}),
	}}
	eng := &Engine{client: script, toolReg: reg}
	res := runQuery(t, eng, QueryParams{
		Messages:  []*types.Message{types.NewUserMessage("loop")},
		Tools:     []tools.Tool{tool},
		Model:     "test-model",
		MaxTokens: 16,
		MaxTurns:  2,
	})
	if res.Reason != "max_turns" {
		t.Fatalf("reason = %q err=%v", res.Reason, res.Error)
	}
	if len(script.bodies()) != 2 {
		t.Fatalf("model requests = %d", len(script.bodies()))
	}
}
