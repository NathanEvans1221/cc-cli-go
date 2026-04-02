package errors

import (
	"fmt"
	"net/http"
)

func NewAPIError(code, message string) *Error {
	return New(ErrorTypeAPI, code, message)
}

func WrapAPIError(err error, code, message string) *Error {
	return Wrap(err, ErrorTypeAPI, code, message)
}

func APIConnectionError(err error) *Error {
	return WrapAPIError(err, "CONNECTION_ERROR",
		"Failed to connect to API").
		WithSuggestion("Check your internet connection and ANTHROPIC_API_KEY environment variable")
}

func APIAuthenticationError(err error) *Error {
	return WrapAPIError(err, "AUTH_ERROR",
		"Authentication failed").
		WithSuggestion("Verify your ANTHROPIC_API_KEY is valid and has not expired")
}

func APIRateLimitError(err error) *Error {
	return WrapAPIError(err, "RATE_LIMIT",
		"API rate limit exceeded").
		WithSuggestion("Wait a moment and try again")
}

func APIInvalidResponseError(err error) *Error {
	return WrapAPIError(err, "INVALID_RESPONSE",
		"Received invalid response from API").
		WithSuggestion("This might be a temporary issue. Try again later")
}

func APITimeoutError(err error) *Error {
	return WrapAPIError(err, "TIMEOUT",
		"API request timed out").
		WithSuggestion("The request took too long. Try a simpler query or check your connection")
}

func APIModelNotFoundError(model string) *Error {
	return NewAPIError("MODEL_NOT_FOUND",
		fmt.Sprintf("Model '%s' not found", model)).
		WithSuggestion("Check the model name or use a different model")
}

func IsAPIStatusError(err error, statusCode int) bool {
	if !IsAPIError(err) {
		return false
	}
	if e, ok := err.(*Error); ok {
		if code, exists := e.Context["status_code"]; exists {
			if sc, ok := code.(int); ok {
				return sc == statusCode
			}
		}
	}
	return false
}

func APIErrorFromStatusCode(statusCode int, err error) *Error {
	switch statusCode {
	case http.StatusUnauthorized:
		return APIAuthenticationError(err)
	case http.StatusTooManyRequests:
		return APIRateLimitError(err)
	case http.StatusRequestTimeout:
		return APITimeoutError(err)
	default:
		return WrapAPIError(err, fmt.Sprintf("HTTP_%d", statusCode),
			fmt.Sprintf("API returned status %d", statusCode))
	}
}
