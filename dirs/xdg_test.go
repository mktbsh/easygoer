package dirs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolve_DefaultsToXDGSpecDefaults(t *testing.T) {
	home := t.TempDir()
	// Clear XDG vars to force defaults
	for _, k := range []string{"XDG_CONFIG_HOME", "XDG_DATA_HOME", "XDG_STATE_HOME", "XDG_CACHE_HOME", "XDG_RUNTIME_DIR"} {
		t.Setenv(k, "")
	}

	p, err := Resolve("mycli", WithHomeDir(home))
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	if p.ConfigHome != filepath.Join(home, ".config") {
		t.Fatalf("ConfigHome = %q", p.ConfigHome)
	}
	if p.DataHome != filepath.Join(home, ".local", "share") {
		t.Fatalf("DataHome = %q", p.DataHome)
	}
	if p.StateHome != filepath.Join(home, ".local", "state") {
		t.Fatalf("StateHome = %q", p.StateHome)
	}
	if p.CacheHome != filepath.Join(home, ".cache") {
		t.Fatalf("CacheHome = %q", p.CacheHome)
	}
	if p.RuntimeDir != "" {
		t.Fatalf("RuntimeDir expected empty, got %q", p.RuntimeDir)
	}

	// app scoped
	if p.ConfigDir != filepath.Join(home, ".config", "mycli") {
		t.Fatalf("ConfigDir = %q", p.ConfigDir)
	}
	if p.PKIDir != filepath.Join(home, ".local", "share", "mycli", "pki") {
		t.Fatalf("PKIDir = %q", p.PKIDir)
	}
	if p.LogDir != filepath.Join(home, ".local", "state", "mycli", "logs") {
		t.Fatalf("LogDir = %q", p.LogDir)
	}
}

func TestResolve_RespectsXDGEnvVars(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "C")
	data := filepath.Join(home, "D")
	state := filepath.Join(home, "S")
	cache := filepath.Join(home, "K")
	runtime := filepath.Join(home, "R")

	t.Setenv("XDG_CONFIG_HOME", cfg)
	t.Setenv("XDG_DATA_HOME", data)
	t.Setenv("XDG_STATE_HOME", state)
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("XDG_RUNTIME_DIR", runtime)

	p, err := Resolve("app", WithHomeDir(home))
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if p.ConfigHome != cfg {
		t.Fatalf("ConfigHome = %q", p.ConfigHome)
	}
	if p.DataHome != data {
		t.Fatalf("DataHome = %q", p.DataHome)
	}
	if p.StateHome != state {
		t.Fatalf("StateHome = %q", p.StateHome)
	}
	if p.CacheHome != cache {
		t.Fatalf("CacheHome = %q", p.CacheHome)
	}
	if p.RuntimeDir != runtime {
		t.Fatalf("RuntimeDir = %q", p.RuntimeDir)
	}
}

func TestResolve_ExpandsTildeAndRelative(t *testing.T) {
	home := t.TempDir()
	// config: ~/c should expand
	t.Setenv("XDG_CONFIG_HOME", "~/c")
	// data: relative should become home-relative
	t.Setenv("XDG_DATA_HOME", "rel/data")
	t.Setenv("XDG_STATE_HOME", "")
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("XDG_RUNTIME_DIR", "")

	p, err := Resolve("app", WithHomeDir(home))
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if p.ConfigHome != filepath.Join(home, "c") {
		t.Fatalf("ConfigHome = %q", p.ConfigHome)
	}
	if p.DataHome != filepath.Join(home, "rel", "data") {
		t.Fatalf("DataHome = %q", p.DataHome)
	}
}

func TestResolve_RequiresAppName(t *testing.T) {
	_, err := Resolve("")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestEnsure_RequiresResolvedPaths(t *testing.T) {
	var p Paths
	if err := p.Ensure(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolve_UsesOSUserHomeDirWhenNoOverride(t *testing.T) {
	// This test is intentionally light: it ensures Resolve does not crash
	// with the real home dir.
	for _, k := range []string{"XDG_CONFIG_HOME", "XDG_DATA_HOME", "XDG_STATE_HOME", "XDG_CACHE_HOME", "XDG_RUNTIME_DIR"} {
		os.Unsetenv(k)
	}
	_, err := Resolve("mycli")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
}
