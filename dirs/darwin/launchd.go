//go:build darwin

package darwin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LaunchAgentsPlistPath returns the per-user LaunchAgent plist path.
//
// macOS expects LaunchAgents here (not XDG):
// ~/Library/LaunchAgents/<label>.plist
func LaunchAgentsPlistPath(label string) (string, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return "", fmt.Errorf("label is required")
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}
	return filepath.Join(home, "Library", "LaunchAgents", label+".plist"), nil
}

// LaunchDaemonsPlistPath returns the system-wide LaunchDaemon plist path.
//
// Note: writing there typically requires admin privileges.
func LaunchDaemonsPlistPath(label string) (string, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return "", fmt.Errorf("label is required")
	}
	return filepath.Join(string(filepath.Separator), "Library", "LaunchDaemons", label+".plist"), nil
}
