package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/liao-eli/cc-cli-go/internal/permission"
)

type Settings struct {
	Permission PermissionSettings `json:"permission"`
	Tools      ToolSettings       `json:"tools"`
	API        APISettings        `json:"api"`
}

type PermissionSettings struct {
	Mode  string           `json:"mode"`
	Rules []PermissionRule `json:"rules"`
}

type PermissionRule struct {
	ToolName string `json:"tool_name"`
	Pattern  string `json:"pattern"`
	Behavior string `json:"behavior"`
}

type ToolSettings struct {
	Enabled  []string `json:"enabled"`
	Disabled []string `json:"disabled"`
}

type APISettings struct {
	Model    string `json:"model"`
	MaxToken int    `json:"max_tokens"`
}

func DefaultSettings() *Settings {
	return &Settings{
		Permission: PermissionSettings{
			Mode: "default",
			Rules: []PermissionRule{
				{ToolName: "Read", Pattern: "*", Behavior: "allow"},
				{ToolName: "Glob", Pattern: "*", Behavior: "allow"},
				{ToolName: "Grep", Pattern: "*", Behavior: "allow"},
			},
		},
		Tools: ToolSettings{
			Enabled:  []string{},
			Disabled: []string{},
		},
		API: APISettings{
			Model:    "claude-3-5-sonnet-20241022",
			MaxToken: 4096,
		},
	}
}

func (s *Settings) ToPermissionRules() []permission.Rule {
	rules := make([]permission.Rule, len(s.Permission.Rules))
	for i, r := range s.Permission.Rules {
		rules[i] = permission.Rule{
			ToolName: r.ToolName,
			Pattern:  r.Pattern,
			Behavior: permission.Behavior(r.Behavior),
		}
	}
	return rules
}

func (s *Settings) GetPermissionMode() permission.Mode {
	return permission.Mode(s.Permission.Mode)
}

func Load() (*Settings, error) {
	settings := DefaultSettings()

	globalSettings, err := loadGlobalSettings()
	if err == nil {
		settings = mergeSettings(settings, globalSettings)
	}

	projectSettings, err := loadProjectSettings()
	if err == nil {
		settings = mergeSettings(settings, projectSettings)
	}

	return settings, nil
}

func loadGlobalSettings() (*Settings, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	globalPath := filepath.Join(homeDir, ".claude", "settings.json")
	return loadSettingsFile(globalPath)
}

func loadProjectSettings() (*Settings, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	projectPath := filepath.Join(cwd, ".claude", "settings.json")
	return loadSettingsFile(projectPath)
}

func loadSettingsFile(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func SaveGlobal(settings *Settings) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".claude")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(configDir, "settings.json")
	return saveSettingsFile(path, settings)
}

func SaveProject(settings *Settings) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	configDir := filepath.Join(cwd, ".claude")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(configDir, "settings.json")
	return saveSettingsFile(path, settings)
}

func saveSettingsFile(path string, settings *Settings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func mergeSettings(base, override *Settings) *Settings {
	if override.Permission.Mode != "" {
		base.Permission.Mode = override.Permission.Mode
	}

	if len(override.Permission.Rules) > 0 {
		base.Permission.Rules = override.Permission.Rules
	}

	if len(override.Tools.Enabled) > 0 {
		base.Tools.Enabled = override.Tools.Enabled
	}

	if len(override.Tools.Disabled) > 0 {
		base.Tools.Disabled = override.Tools.Disabled
	}

	if override.API.Model != "" {
		base.API.Model = override.API.Model
	}

	if override.API.MaxToken > 0 {
		base.API.MaxToken = override.API.MaxToken
	}

	return base
}
