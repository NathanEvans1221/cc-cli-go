package config

import (
	"fmt"

	"github.com/liao-eli/cc-cli-go/internal/permission"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (s *Settings) Validate() []ValidationError {
	var errors []ValidationError

	errors = append(errors, s.validatePermissionMode()...)
	errors = append(errors, s.validatePermissionRules()...)
	errors = append(errors, s.validateAPISettings()...)

	return errors
}

func (s *Settings) validatePermissionMode() []ValidationError {
	var errors []ValidationError

	validModes := map[string]bool{
		string(permission.ModeDefault): true,
		string(permission.ModeAccept):  true,
		string(permission.ModePlan):    true,
		string(permission.ModeAuto):    true,
	}

	if s.Permission.Mode != "" && !validModes[s.Permission.Mode] {
		errors = append(errors, ValidationError{
			Field:   "permission.mode",
			Message: fmt.Sprintf("invalid mode '%s', must be one of: default, accept, plan, auto", s.Permission.Mode),
		})
	}

	return errors
}

func (s *Settings) validatePermissionRules() []ValidationError {
	var errors []ValidationError

	validBehaviors := map[string]bool{
		string(permission.BehaviorAllow): true,
		string(permission.BehaviorDeny):  true,
		string(permission.BehaviorAsk):   true,
	}

	for i, rule := range s.Permission.Rules {
		if rule.ToolName == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("permission.rules[%d].tool_name", i),
				Message: "tool_name is required",
			})
		}

		if rule.Pattern == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("permission.rules[%d].pattern", i),
				Message: "pattern is required",
			})
		}

		if rule.Behavior == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("permission.rules[%d].behavior", i),
				Message: "behavior is required",
			})
		} else if !validBehaviors[rule.Behavior] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("permission.rules[%d].behavior", i),
				Message: fmt.Sprintf("invalid behavior '%s', must be one of: allow, deny, ask", rule.Behavior),
			})
		}
	}

	return errors
}

func (s *Settings) validateAPISettings() []ValidationError {
	var errors []ValidationError

	if s.API.Model != "" && len(s.API.Model) < 5 {
		errors = append(errors, ValidationError{
			Field:   "api.model",
			Message: "model name too short",
		})
	}

	if s.API.MaxToken < 0 {
		errors = append(errors, ValidationError{
			Field:   "api.max_tokens",
			Message: "max_tokens must be non-negative",
		})
	}

	if s.API.MaxToken > 100000 {
		errors = append(errors, ValidationError{
			Field:   "api.max_tokens",
			Message: "max_tokens exceeds maximum limit (100000)",
		})
	}

	return errors
}

func (s *Settings) IsValid() bool {
	return len(s.Validate()) == 0
}

func FormatValidationErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	result := "Validation errors:\n"
	for _, err := range errors {
		result += fmt.Sprintf("  - %s\n", err.Error())
	}
	return result
}
