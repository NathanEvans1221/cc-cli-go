package errors

import (
	"fmt"
)

func NewToolError(toolName, code, message string) *Error {
	return New(ErrorTypeTool, code, message).
		WithContext("tool_name", toolName)
}

func WrapToolError(err error, toolName, code, message string) *Error {
	return Wrap(err, ErrorTypeTool, code, message).
		WithContext("tool_name", toolName)
}

func ToolNotFoundError(toolName string) *Error {
	return NewToolError(toolName, "TOOL_NOT_FOUND",
		fmt.Sprintf("Tool '%s' not found", toolName)).
		WithSuggestion("Check the tool name or register the tool")
}

func ToolInputValidationError(toolName, field, reason string) *Error {
	return NewToolError(toolName, "INVALID_INPUT",
		fmt.Sprintf("Invalid input for field '%s': %s", field, reason)).
		WithSuggestion("Check the tool's input schema and provide valid parameters")
}

func ToolExecutionError(toolName string, err error) *Error {
	return WrapToolError(err, toolName, "EXECUTION_ERROR",
		fmt.Sprintf("Tool '%s' execution failed", toolName)).
		WithSuggestion("Check the tool input and try again")
}

func ToolPermissionDeniedError(toolName string, reason string) *Error {
	return NewToolError(toolName, "PERMISSION_DENIED",
		fmt.Sprintf("Permission denied: %s", reason)).
		WithSuggestion("This tool requires permission. Use permission mode 'accept' or grant specific permissions")
}

func ToolTimeoutError(toolName string) *Error {
	return NewToolError(toolName, "TIMEOUT",
		"Tool execution timed out").
		WithSuggestion("The tool took too long. Try simplifying the operation or increasing timeout")
}

func ToolFileNotFoundError(toolName, filePath string) *Error {
	return NewToolError(toolName, "FILE_NOT_FOUND",
		fmt.Sprintf("File not found: %s", filePath)).
		WithSuggestion("Check the file path and ensure the file exists")
}

func ToolInvalidPathError(toolName, filePath string) *Error {
	return NewToolError(toolName, "INVALID_PATH",
		fmt.Sprintf("Invalid file path: %s", filePath)).
		WithSuggestion("Provide an absolute path or a valid relative path")
}

func ToolCommandError(toolName, command string, err error) *Error {
	return WrapToolError(err, toolName, "COMMAND_ERROR",
		fmt.Sprintf("Command failed: %s", command)).
		WithSuggestion("Check the command syntax and ensure all dependencies are installed")
}
