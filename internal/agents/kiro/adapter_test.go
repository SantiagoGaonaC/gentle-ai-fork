package kiro

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

func TestAdapter_Agent(t *testing.T) {
	adapter := NewAdapter()
	if got := adapter.Agent(); got != model.AgentKiroIDE {
		t.Errorf("Agent() = %q, want %q", got, model.AgentKiroIDE)
	}
}

func TestAdapter_Tier(t *testing.T) {
	adapter := NewAdapter()
	if got := adapter.Tier(); got != model.TierFull {
		t.Errorf("Tier() = %q, want %q", got, model.TierFull)
	}
}

func TestAdapter_Detect_BinaryNotFound(t *testing.T) {
	adapter := &Adapter{
		lookPath: func(string) (string, error) {
			return "", &mockLookPathError{}
		},
	}

	installed, _, _, _, err := adapter.Detect(nil, "")
	if installed {
		t.Error("Detect() should return false when binary not found")
	}
	if err != nil {
		t.Errorf("Detect() should not return error when binary not found, got %v", err)
	}
}

func TestAdapter_Detect_BinaryFound(t *testing.T) {
	adapter := &Adapter{
		lookPath: func(string) (string, error) {
			return "/usr/local/bin/kiro", nil
		},
	}

	installed, binaryPath, _, configFound, err := adapter.Detect(nil, "")
	if !installed {
		t.Error("Detect() should return true when binary found")
	}
	if binaryPath != "/usr/local/bin/kiro" {
		t.Errorf("Detect() binaryPath = %q, want %q", binaryPath, "/usr/local/bin/kiro")
	}
	if !configFound {
		t.Error("Detect() configFound should be true")
	}
	if err != nil {
		t.Errorf("Detect() should not return error, got %v", err)
	}
}

func TestAdapter_GlobalConfigDir(t *testing.T) {
	adapter := NewAdapter()
	homeDir := "/home/user"
	got := adapter.GlobalConfigDir(homeDir)

	// Verify path ends with expected structure based on OS
	switch runtime.GOOS {
	case "darwin":
		if !contains(got, "Library", "Application Support", "Kiro", "User") {
			t.Errorf("macOS: GlobalConfigDir() = %q, missing expected path components", got)
		}
	case "windows":
		if !contains(got, "kiro", "User") {
			t.Errorf("Windows: GlobalConfigDir() = %q, missing expected path components", got)
		}
	default: // linux
		if !contains(got, "kiro", "user") {
			t.Errorf("Linux: GlobalConfigDir() = %q, missing expected path components", got)
		}
	}
}

func TestAdapter_SystemPromptDir(t *testing.T) {
	adapter := NewAdapter()
	homeDir := "/home/user"
	got := adapter.SystemPromptDir(homeDir)
	expected := filepath.Join(homeDir, ".kiro", "steering")

	if got != expected {
		t.Errorf("SystemPromptDir() = %q, want %q", got, expected)
	}
}

func TestAdapter_SystemPromptFile(t *testing.T) {
	adapter := NewAdapter()
	homeDir := "/home/user"
	expected := filepath.Join(homeDir, ".kiro", "steering", "gentle-ai.md")

	got := adapter.SystemPromptFile(homeDir)
	if got != expected {
		t.Errorf("SystemPromptFile() = %q, want %q", got, expected)
	}
}

func TestAdapter_SkillsDir(t *testing.T) {
	adapter := NewAdapter()
	homeDir := "/home/user"
	expected := filepath.Join(homeDir, ".kiro", "skills")

	got := adapter.SkillsDir(homeDir)
	if got != expected {
		t.Errorf("SkillsDir() = %q, want %q", got, expected)
	}

	// Verify path is independent from GlobalConfigDir (must not contain AppData or platform config dir).
	globalConfigDir := adapter.GlobalConfigDir(homeDir)
	if got == filepath.Join(globalConfigDir, "skills") {
		t.Errorf("SkillsDir() must be independent from GlobalConfigDir(); got %q which matches GlobalConfigDir/skills", got)
	}
}

func TestAdapter_SettingsPath(t *testing.T) {
	adapter := NewAdapter()
	homeDir := "/home/user"
	configDir := adapter.GlobalConfigDir(homeDir)
	expected := filepath.Join(configDir, "settings.json")

	got := adapter.SettingsPath(homeDir)
	if got != expected {
		t.Errorf("SettingsPath() = %q, want %q", got, expected)
	}
}

func TestAdapter_MCPConfigPath(t *testing.T) {
	adapter := NewAdapter()
	homeDir := "/home/user"
	// Kiro reads MCP from ~/.kiro/settings/mcp.json, not from the app config dir.
	expected := filepath.Join(homeDir, ".kiro", "settings", "mcp.json")

	got := adapter.MCPConfigPath(homeDir, "")
	if got != expected {
		t.Errorf("MCPConfigPath() = %q, want %q", got, expected)
	}
}

func TestAdapter_SystemPromptStrategy(t *testing.T) {
	adapter := NewAdapter()
	expected := model.StrategySteeringFile

	got := adapter.SystemPromptStrategy()
	if got != expected {
		t.Errorf("SystemPromptStrategy() = %v, want %v", got, expected)
	}
}

func TestAdapter_SupportsSubAgents(t *testing.T) {
	adapter := NewAdapter()
	if !adapter.SupportsSubAgents() {
		t.Error("SupportsSubAgents() should return true")
	}
}

func TestAdapter_SubAgentsDir(t *testing.T) {
	adapter := NewAdapter()
	homeDir := "/home/user"
	expected := filepath.Join(homeDir, ".kiro", "agents")
	if got := adapter.SubAgentsDir(homeDir); got != expected {
		t.Errorf("SubAgentsDir() = %q, want %q", got, expected)
	}
}

func TestAdapter_EmbeddedSubAgentsDir(t *testing.T) {
	adapter := NewAdapter()
	if got := adapter.EmbeddedSubAgentsDir(); got != "kiro/agents" {
		t.Errorf("EmbeddedSubAgentsDir() = %q, want %q", got, "kiro/agents")
	}
}

func TestAdapter_KiroModelID(t *testing.T) {
	adapter := NewAdapter()
	tests := []struct {
		alias model.ClaudeModelAlias
		want  string
	}{
		{model.ClaudeModelOpus, "claude-opus-4.6"},
		{model.ClaudeModelSonnet, "claude-sonnet-4.6"},
		{model.ClaudeModelHaiku, "claude-haiku-4.5"},
		{"unknown", "claude-sonnet-4.6"},
	}
	for _, tt := range tests {
		if got := adapter.KiroModelID(tt.alias); got != tt.want {
			t.Errorf("KiroModelID(%q) = %v, want %v", tt.alias, got, tt.want)
		}
	}
}

func TestAdapter_MCPStrategy(t *testing.T) {
	adapter := NewAdapter()
	expected := model.StrategyMCPConfigFile

	got := adapter.MCPStrategy()
	if got != expected {
		t.Errorf("MCPStrategy() = %q, want %q", got, expected)
	}
}

func TestAdapter_InstallCommand_macOS(t *testing.T) {
	adapter := NewAdapter()
	profile := system.PlatformProfile{OS: "darwin"}

	_, err := adapter.InstallCommand(profile)
	if err == nil {
		t.Error("InstallCommand() should return error (auto-install not supported)")
	}
	if _, ok := err.(AgentNotInstallableError); !ok {
		t.Errorf("InstallCommand() expected AgentNotInstallableError, got %T", err)
	}
}

func TestAdapter_InstallCommand_Linux(t *testing.T) {
	adapter := NewAdapter()
	profile := system.PlatformProfile{OS: "linux"}

	_, err := adapter.InstallCommand(profile)
	if err == nil {
		t.Error("InstallCommand() should return error (auto-install not supported)")
	}
	if _, ok := err.(AgentNotInstallableError); !ok {
		t.Errorf("InstallCommand() expected AgentNotInstallableError, got %T", err)
	}
}

func TestAdapter_InstallCommand_Windows(t *testing.T) {
	adapter := NewAdapter()
	profile := system.PlatformProfile{OS: "windows"}

	_, err := adapter.InstallCommand(profile)
	if err == nil {
		t.Error("InstallCommand() should return error (auto-install not supported)")
	}
	if _, ok := err.(AgentNotInstallableError); !ok {
		t.Errorf("InstallCommand() expected AgentNotInstallableError, got %T", err)
	}
}

func TestAdapter_SupportsFeatures(t *testing.T) {
	adapter := NewAdapter()

	tests := []struct {
		name     string
		fn       func() bool
		expected bool
	}{
		{"SupportsSkills", adapter.SupportsSkills, true},
		{"SupportsSystemPrompt", adapter.SupportsSystemPrompt, true},
		{"SupportsMCP", adapter.SupportsMCP, true},
		{"SupportsOutputStyles", adapter.SupportsOutputStyles, false},
		{"SupportsSlashCommands", adapter.SupportsSlashCommands, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn()
			if got != tt.expected {
				t.Errorf("%s() = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

// mockLookPathError is a mock error for testing
type mockLookPathError struct{}

func (e *mockLookPathError) Error() string { return "executable not found" }
func (e *mockLookPathError) Unwrap() error { return nil }

// contains checks if a path contains all given components as substrings
func contains(path string, components ...string) bool {
	for _, comp := range components {
		if !stringContains(path, comp) {
			return false
		}
	}
	return true
}

// stringContains is a simple substring check
func stringContains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
