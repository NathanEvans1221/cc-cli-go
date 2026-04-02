// Package permission provides a permission management system for tool execution.
// It allows fine-grained control over which tools can be executed and under what conditions.
package permission

// Mode defines the overall permission checking strategy.
type Mode string

const (
	// ModeDefault asks for permission on potentially dangerous operations
	ModeDefault Mode = "default"
	// ModeAccept automatically allows all operations
	ModeAccept Mode = "accept"
	// ModePlan requires approval for all operations
	ModePlan Mode = "plan"
	// ModeAuto allows safe operations, asks for dangerous ones
	ModeAuto Mode = "auto"
)

// Behavior defines the action to take for a permission request.
type Behavior string

const (
	// BehaviorAllow grants permission without asking
	BehaviorAllow Behavior = "allow"
	// BehaviorDeny refuses permission without asking
	BehaviorDeny Behavior = "deny"
	// BehaviorAsk prompts the user for a decision
	BehaviorAsk Behavior = "ask"
)

// Rule defines a permission rule for a specific tool and input pattern.
type Rule struct {
	// ToolName specifies which tool this rule applies to
	ToolName string
	// Pattern matches against tool input (supports wildcards)
	Pattern string
	// Behavior defines the action when the rule matches
	Behavior Behavior
}

// Decision represents the result of a permission check.
type Decision struct {
	// Behavior indicates what action to take
	Behavior Behavior
	// Reason provides context for the decision
	Reason string
}

// Checker manages permission rules and modes for tool execution.
type Checker struct {
	mode  Mode
	rules []Rule
}

// NewChecker creates a new permission checker with the specified mode.
func NewChecker(mode Mode) *Checker {
	return &Checker{
		mode:  mode,
		rules: getDefaultRules(),
	}
}

// getDefaultRules returns the default permission rules.
// By default, read-only tools (Read, Glob, Grep) are allowed.
func getDefaultRules() []Rule {
	return []Rule{
		{ToolName: "Read", Pattern: "*", Behavior: BehaviorAllow},
		{ToolName: "Glob", Pattern: "*", Behavior: BehaviorAllow},
		{ToolName: "Grep", Pattern: "*", Behavior: BehaviorAllow},
	}
}

func (c *Checker) SetRules(rules []Rule) {
	c.rules = rules
}

func (c *Checker) Check(toolName string, input map[string]interface{}) *Decision {
	if c.mode == ModeAccept {
		return &Decision{Behavior: BehaviorAllow, Reason: "accept mode"}
	}

	if c.mode == ModeAuto {
		return c.checkAutoMode(toolName, input)
	}

	for _, rule := range c.rules {
		if rule.ToolName != toolName {
			continue
		}

		if matchPattern(rule.Pattern, input) {
			return &Decision{Behavior: rule.Behavior, Reason: "rule matched"}
		}
	}

	if isDangerousCommand(toolName, input) {
		return &Decision{Behavior: BehaviorAsk, Reason: "dangerous command detected"}
	}

	return &Decision{Behavior: BehaviorAsk, Reason: "no matching rule"}
}

func (c *Checker) checkAutoMode(toolName string, input map[string]interface{}) *Decision {
	if isDangerousCommand(toolName, input) {
		return &Decision{Behavior: BehaviorAsk, Reason: "dangerous command in auto mode"}
	}

	return &Decision{Behavior: BehaviorAllow, Reason: "auto mode"}
}

func matchPattern(pattern string, input map[string]interface{}) bool {
	if pattern == "*" {
		return true
	}

	for _, value := range input {
		if strValue, ok := value.(string); ok {
			if matchStringPattern(pattern, strValue) {
				return true
			}
		}
	}

	return false
}

func matchStringPattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}

	return value == pattern || containsPattern(pattern, value)
}

func containsPattern(pattern, value string) bool {
	return len(pattern) > 0 && len(value) >= len(pattern) && value[:len(pattern)] == pattern
}
