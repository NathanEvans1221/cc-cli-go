package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
)

func TestDefaultSettings(t *testing.T) {
	settings := DefaultSettings()

	if settings == nil {
		t.Fatal("expected non-nil settings")
	}

	testutil.AssertEqual(t, "default", settings.Permission.Mode)
	testutil.AssertEqual(t, "claude-3-5-sonnet-20241022", settings.API.Model)
	testutil.AssertEqual(t, 4096, settings.API.MaxToken)
}

func TestDefaultSettings_HasDefaultRules(t *testing.T) {
	settings := DefaultSettings()

	if len(settings.Permission.Rules) == 0 {
		t.Error("expected default permission rules")
	}
}

func TestSettings_Validate_ValidSettings(t *testing.T) {
	settings := DefaultSettings()

	errors := settings.Validate()

	if len(errors) > 0 {
		t.Errorf("expected no validation errors, got %d", len(errors))
	}
}

func TestSettings_Validate_InvalidPermissionMode(t *testing.T) {
	settings := DefaultSettings()
	settings.Permission.Mode = "invalid_mode"

	errors := settings.Validate()

	if len(errors) == 0 {
		t.Error("expected validation error for invalid mode")
	}

	found := false
	for _, err := range errors {
		if err.Field == "permission.mode" {
			found = true
		}
	}

	if !found {
		t.Error("expected validation error for permission.mode")
	}
}

func TestSettings_Validate_InvalidBehavior(t *testing.T) {
	settings := DefaultSettings()
	settings.Permission.Rules = []PermissionRule{
		{ToolName: "Bash", Pattern: "*", Behavior: "invalid"},
	}

	errors := settings.Validate()

	if len(errors) == 0 {
		t.Error("expected validation error for invalid behavior")
	}
}

func TestSettings_Validate_MissingToolName(t *testing.T) {
	settings := DefaultSettings()
	settings.Permission.Rules = []PermissionRule{
		{ToolName: "", Pattern: "*", Behavior: "allow"},
	}

	errors := settings.Validate()

	if len(errors) == 0 {
		t.Error("expected validation error for missing tool_name")
	}
}

func TestSettings_Validate_InvalidMaxTokens(t *testing.T) {
	settings := DefaultSettings()
	settings.API.MaxToken = -100

	errors := settings.Validate()

	if len(errors) == 0 {
		t.Error("expected validation error for negative max_tokens")
	}
}

func TestSettings_Validate_MaxTokensTooLarge(t *testing.T) {
	settings := DefaultSettings()
	settings.API.MaxToken = 200000

	errors := settings.Validate()

	if len(errors) == 0 {
		t.Error("expected validation error for max_tokens exceeding limit")
	}
}

func TestSettings_IsValid(t *testing.T) {
	settings := DefaultSettings()

	if !settings.IsValid() {
		t.Error("expected default settings to be valid")
	}
}

func TestSettings_ToPermissionRules(t *testing.T) {
	settings := &Settings{
		Permission: PermissionSettings{
			Rules: []PermissionRule{
				{ToolName: "Read", Pattern: "*", Behavior: "allow"},
				{ToolName: "Bash", Pattern: "ls", Behavior: "ask"},
			},
		},
	}

	rules := settings.ToPermissionRules()

	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}

	if rules[0].ToolName != "Read" {
		t.Errorf("expected Read, got %s", rules[0].ToolName)
	}
}

func TestMergeSettings_OverrideMode(t *testing.T) {
	base := DefaultSettings()
	override := &Settings{
		Permission: PermissionSettings{
			Mode: "accept",
		},
	}

	merged := mergeSettings(base, override)

	if merged.Permission.Mode != "accept" {
		t.Errorf("expected accept mode, got %s", merged.Permission.Mode)
	}
}

func TestMergeSettings_OverrideRules(t *testing.T) {
	base := DefaultSettings()
	override := &Settings{
		Permission: PermissionSettings{
			Rules: []PermissionRule{
				{ToolName: "Custom", Pattern: "*", Behavior: "deny"},
			},
		},
	}

	merged := mergeSettings(base, override)

	if len(merged.Permission.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(merged.Permission.Rules))
	}
}

func TestMergeSettings_OverrideAPI(t *testing.T) {
	base := DefaultSettings()
	override := &Settings{
		API: APISettings{
			Model:    "claude-3-opus",
			MaxToken: 8192,
		},
	}

	merged := mergeSettings(base, override)

	if merged.API.Model != "claude-3-opus" {
		t.Errorf("expected claude-3-opus, got %s", merged.API.Model)
	}

	if merged.API.MaxToken != 8192 {
		t.Errorf("expected 8192, got %d", merged.API.MaxToken)
	}
}

func TestLoadSettingsFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "settings.json")

	content := `{
		"permission": {
			"mode": "accept"
		},
		"api": {
			"model": "claude-3-opus",
			"max_tokens": 8192
		}
	}`

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	settings, err := loadSettingsFile(configPath)
	testutil.AssertNoError(t, err)

	if settings.Permission.Mode != "accept" {
		t.Errorf("expected accept, got %s", settings.Permission.Mode)
	}

	if settings.API.Model != "claude-3-opus" {
		t.Errorf("expected claude-3-opus, got %s", settings.API.Model)
	}
}

func TestLoadSettingsFile_NotFound(t *testing.T) {
	_, err := loadSettingsFile("/nonexistent/settings.json")

	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestSaveSettingsFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "settings.json")

	settings := &Settings{
		Permission: PermissionSettings{
			Mode: "plan",
		},
		API: APISettings{
			Model:    "claude-3-sonnet",
			MaxToken: 4096,
		},
	}

	err := saveSettingsFile(configPath, settings)
	testutil.AssertNoError(t, err)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected settings file to be created")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "settings.json")

	original := &Settings{
		Permission: PermissionSettings{
			Mode: "auto",
			Rules: []PermissionRule{
				{ToolName: "Bash", Pattern: "git*", Behavior: "allow"},
			},
		},
		API: APISettings{
			Model:    "claude-3-opus",
			MaxToken: 8192,
		},
	}

	err := saveSettingsFile(configPath, original)
	testutil.AssertNoError(t, err)

	loaded, err := loadSettingsFile(configPath)
	testutil.AssertNoError(t, err)

	testutil.AssertEqual(t, original.Permission.Mode, loaded.Permission.Mode)
	testutil.AssertEqual(t, original.API.Model, loaded.API.Model)
	testutil.AssertEqual(t, original.API.MaxToken, loaded.API.MaxToken)
}

func TestLoad_GlobalSettings(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	configDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := `{
		"permission": {
			"mode": "accept"
		}
	}`

	configPath := filepath.Join(configDir, "settings.json")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)

	settings, err := Load()
	testutil.AssertNoError(t, err)

	if settings.Permission.Mode != "accept" {
		t.Errorf("expected accept, got %s", settings.Permission.Mode)
	}
}

func TestFormatValidationErrors(t *testing.T) {
	errors := []ValidationError{
		{Field: "test", Message: "error message"},
	}

	result := FormatValidationErrors(errors)

	if result == "" {
		t.Error("expected non-empty formatted errors")
	}

	testutil.AssertContains(t, result, "test")
	testutil.AssertContains(t, result, "error message")
}

func TestFormatValidationErrors_Empty(t *testing.T) {
	result := FormatValidationErrors([]ValidationError{})

	if result != "" {
		t.Error("expected empty string for no errors")
	}
}
