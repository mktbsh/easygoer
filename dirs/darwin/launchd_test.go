//go:build darwin

package darwin

import (
	"path/filepath"
	"testing"
)

func TestLaunchAgentsPlistPath(t *testing.T) {
	home := t.TempDir()
	// os.UserHomeDir() reads HOME on darwin; set it.
	t.Setenv("HOME", home)

	p, err := LaunchAgentsPlistPath("com.example.mycli")
	if err != nil {
		t.Fatalf("LaunchAgentsPlistPath: %v", err)
	}
	want := filepath.Join(home, "Library", "LaunchAgents", "com.example.mycli.plist")
	if p != want {
		t.Fatalf("got %q, want %q", p, want)
	}
}

func TestLaunchDaemonsPlistPath(t *testing.T) {
	p, err := LaunchDaemonsPlistPath("com.example.mycli")
	if err != nil {
		t.Fatalf("LaunchDaemonsPlistPath: %v", err)
	}
	want := filepath.Join(string(filepath.Separator), "Library", "LaunchDaemons", "com.example.mycli.plist")
	if p != want {
		t.Fatalf("got %q, want %q", p, want)
	}
}

func TestLaunchAgentsPlistPath_RequiresLabel(t *testing.T) {
	_, err := LaunchAgentsPlistPath(" ")
	if err == nil {
		t.Fatalf("expected error")
	}
}
