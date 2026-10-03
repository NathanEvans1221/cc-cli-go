package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorFormattingAndPredicates(t *testing.T) {
	base := errors.New("dial failed")
	err := Wrap(base, ErrorTypeAPI, "CONNECTION_ERROR", "cannot reach API").
		WithContext("status_code", 503).
		WithSuggestion("retry later")

	if !strings.Contains(err.Error(), "CONNECTION_ERROR") || !strings.Contains(err.Error(), "dial failed") {
		t.Fatalf("error string %s", err.Error())
	}
	if !errors.Is(err, base) {
		t.Fatal("unwrap did not expose the cause")
	}
	if !strings.Contains(err.UserMessage(), "retry later") {
		t.Fatalf("user message %s", err.UserMessage())
	}
	if !IsAPIError(err) || IsToolError(err) {
		t.Fatal("type predicates")
	}
	if GetErrorType(err) != ErrorTypeAPI || GetErrorType(nil) != ErrorTypeInternal || GetErrorType(base) != ErrorTypeInternal {
		t.Fatal("GetErrorType")
	}
	if IsAPIError(nil) || IsToolError(nil) || IsAPIError(base) {
		t.Fatal("nil and plain errors are not typed")
	}

	plain := New(ErrorTypeSession, "MISSING", "")
	if !strings.Contains(plain.Error(), "MISSING") || strings.Contains(plain.Error(), ":") {
		t.Fatalf("plain %s", plain.Error())
	}
	if plain.UserMessage() != "" {
		t.Fatalf("empty user message %q", plain.UserMessage())
	}
}

func TestAPIAndToolConstructors(t *testing.T) {
	cause := errors.New("boom")
	apiErr := APIErrorFromStatusCode(401, cause)
	if apiErr.Code != "AUTH_ERROR" || !IsAPIStatusError(apiErr.WithContext("status_code", 401), 401) {
		t.Fatalf("auth %+v", apiErr)
	}
	if IsAPIStatusError(apiErr, 500) || IsAPIStatusError(cause, 401) {
		t.Fatal("status mismatch")
	}
	limited := APIErrorFromStatusCode(429, cause)
	timeout := APIErrorFromStatusCode(408, cause)
	other := APIErrorFromStatusCode(500, cause)
	if limited.Code != "RATE_LIMIT" || timeout.Code != "TIMEOUT" || !strings.Contains(other.Code, "HTTP_500") {
		t.Fatalf("%s %s %s", limited.Code, timeout.Code, other.Code)
	}
	if APIConnectionError(cause).Suggestion == "" || APIInvalidResponseError(cause).Code != "INVALID_RESPONSE" {
		t.Fatal("connection constructors")
	}
	if !strings.Contains(APIModelNotFoundError("missing").Message, "missing") {
		t.Fatal("model error")
	}

	tool := ToolNotFoundError("Nope")
	if !IsToolError(tool) || tool.Context["tool_name"] != "Nope" {
		t.Fatalf("%+v", tool)
	}
	cases := []*Error{
		ToolInputValidationError("Read", "path", "empty"),
		ToolExecutionError("Read", cause),
		ToolPermissionDeniedError("Bash", "denied"),
		ToolTimeoutError("Bash"),
		ToolFileNotFoundError("Read", "a.txt"),
		ToolInvalidPathError("Write", "../x"),
		ToolCommandError("Bash", "ls", cause),
	}
	for _, item := range cases {
		if item.Error() == "" || item.Suggestion == "" {
			t.Fatalf("incomplete tool error %+v", item)
		}
	}
}
