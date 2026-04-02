package errors

import (
	"fmt"
	"strings"
)

type ErrorType string

const (
	ErrorTypeAPI        ErrorType = "api"
	ErrorTypeTool       ErrorType = "tool"
	ErrorTypePermission ErrorType = "permission"
	ErrorTypeConfig     ErrorType = "config"
	ErrorTypeSession    ErrorType = "session"
	ErrorTypeInternal   ErrorType = "internal"
)

type Error struct {
	Type       ErrorType
	Code       string
	Message    string
	Cause      error
	Context    map[string]interface{}
	Suggestion string
}

func (e *Error) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s] %s", e.Type, e.Code))

	if e.Message != "" {
		sb.WriteString(fmt.Sprintf(": %s", e.Message))
	}

	if e.Cause != nil {
		sb.WriteString(fmt.Sprintf(" (caused by: %v)", e.Cause))
	}

	return sb.String()
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(typ ErrorType, code, message string) *Error {
	return &Error{
		Type:    typ,
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

func Wrap(err error, typ ErrorType, code, message string) *Error {
	return &Error{
		Type:    typ,
		Code:    code,
		Message: message,
		Cause:   err,
		Context: make(map[string]interface{}),
	}
}

func (e *Error) WithContext(key string, value interface{}) *Error {
	e.Context[key] = value
	return e
}

func (e *Error) WithSuggestion(suggestion string) *Error {
	e.Suggestion = suggestion
	return e
}

func (e *Error) UserMessage() string {
	var sb strings.Builder

	sb.WriteString(e.Message)

	if e.Suggestion != "" {
		sb.WriteString(fmt.Sprintf("\n\nSuggestion: %s", e.Suggestion))
	}

	return sb.String()
}

func IsAPIError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeAPI
	}
	return false
}

func IsToolError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeTool
	}
	return false
}

func GetErrorType(err error) ErrorType {
	if err == nil {
		return ErrorTypeInternal
	}
	if e, ok := err.(*Error); ok {
		return e.Type
	}
	return ErrorTypeInternal
}
