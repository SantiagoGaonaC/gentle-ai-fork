package kiro

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gentleman-programming/gentle-ai/internal/model"
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

type Adapter struct {
	lookPath func(string) (string, error)
}

func NewAdapter() *Adapter {
	return &Adapter{
		lookPath: exec.LookPath,
	}
}

// --- Identity ---

func (a *Adapter) Agent() model.AgentID {
	return model.AgentKiroIDE
}

func (a *Adapter) Tier() model.SupportTier {
	return model.TierFull
}

// --- Detection ---

func (a *Adapter) Detect(_ context.Context, _ string) (bool, string, string, bool, error) {
	// Kiro IDE is a VS Code fork available as a desktop application.
	// Official website: https://kiro.dev/
	// Detect by the "kiro" binary on PATH.
	binaryPath, err := a.lookPath("kiro")
	if err != nil {
		return false, "", "", false, nil
	}

	return true, binaryPath, "", true, nil
}

// --- Installation ---

func (a *Adapter) SupportsAutoInstall() bool {
	return false // Kiro IDE is a desktop app, installed via official downloads or package managers
}

func (a *Adapter) InstallCommand(_ system.PlatformProfile) ([][]string, error) {
	return nil, AgentNotInstallableError{Agent: model.AgentKiroIDE}
}

// --- Config paths ---
// Kiro IDE (VS Code fork) uses similar directory structure to VS Code:
//   - macOS: ~/Library/Application Support/Kiro/User/
//   - Linux: ~/.config/kiro/user/ (respects XDG_CONFIG_HOME)
//   - Windows: %APPDATA%/kiro/User/
// System prompts are loaded from .instructions.md files in the prompts folder

func (a *Adapter) GlobalConfigDir(homeDir string) string {
	return a.kiroConfigDir(homeDir)
}

func (a *Adapter) SystemPromptDir(homeDir string) string {
	return filepath.Join(a.kiroConfigDir(homeDir), "prompts")
}

func (a *Adapter) SystemPromptFile(homeDir string) string {
	return filepath.Join(a.SystemPromptDir(homeDir), "gentle-ai.instructions.md")
}

func (a *Adapter) SkillsDir(homeDir string) string {
	// Skills stored in ~/.config/kiro/user/skills (or equivalent per OS)
	return filepath.Join(a.GlobalConfigDir(homeDir), "skills")
}

func (a *Adapter) SettingsPath(homeDir string) string {
	return filepath.Join(a.kiroConfigDir(homeDir), "settings.json")
}

// --- Config strategies ---

func (a *Adapter) SystemPromptStrategy() model.SystemPromptStrategy {
	return model.StrategyInstructionsFile
}

func (a *Adapter) MCPStrategy() model.MCPStrategy {
	return model.StrategyMCPConfigFile
}

// --- MCP ---

// MCPConfigPath returns the user-level MCP config file.
// Kiro reads MCP configuration from ~/.kiro/settings/mcp.json (user level)
// or .kiro/settings/mcp.json (workspace level). This is separate from the
// app config dir (%APPDATA%/kiro/User on Windows) used for settings and prompts.
func (a *Adapter) MCPConfigPath(homeDir string, _ string) string {
	return filepath.Join(homeDir, ".kiro", "settings", "mcp.json")
}

func (a *Adapter) kiroConfigDir(homeDir string) string {
	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Application Support/Kiro/User/
		return filepath.Join(homeDir, "Library", "Application Support", "Kiro", "User")
	case "windows":
		// Windows: %APPDATA%/kiro/User/
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		return filepath.Join(appData, "kiro", "User")
	default:
		// Linux and others: ~/.config/kiro/user (respects XDG_CONFIG_HOME)
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			xdgConfigHome = filepath.Join(homeDir, ".config")
		}
		return filepath.Join(xdgConfigHome, "kiro", "user")
	}
}

// --- Optional capabilities ---

func (a *Adapter) SupportsOutputStyles() bool {
	return false // Kiro IDE output style support not documented
}

func (a *Adapter) OutputStyleDir(_ string) string {
	return ""
}

func (a *Adapter) SupportsSlashCommands() bool {
	return false // Would need to verify if Kiro IDE has slash command support
}

func (a *Adapter) CommandsDir(_ string) string {
	return ""
}

func (a *Adapter) SupportsSkills() bool {
	return true
}

func (a *Adapter) SupportsSystemPrompt() bool {
	return true
}

func (a *Adapter) SupportsMCP() bool {
	return true
}
