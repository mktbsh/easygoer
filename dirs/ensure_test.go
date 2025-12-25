package dirs

import (
	"os"
	"runtime"
	"testing"
)

func TestEnsure_CreatesDirectories(t *testing.T) {
	home := t.TempDir()
	for _, k := range []string{"XDG_CONFIG_HOME", "XDG_DATA_HOME", "XDG_STATE_HOME", "XDG_CACHE_HOME", "XDG_RUNTIME_DIR"} {
		t.Setenv(k, "")
	}

	p, err := Resolve("mycli", WithHomeDir(home))
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if err := p.Ensure(); err != nil {
		t.Fatalf("Ensure: %v", err)
	}

	for _, dir := range []string{p.ConfigDir, p.DataDir, p.StateDir, p.CacheDir, p.RunDir, p.LogDir, p.PKIDir, p.ServiceDir} {
		st, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("stat %s: %v", dir, err)
		}
		if !st.IsDir() {
			t.Fatalf("%s is not a directory", dir)
		}

		// Permission checks are only meaningful on unix-like systems.
		if runtime.GOOS != "windows" {
			mode := st.Mode().Perm()
			if mode != 0o700 {
				// Umask is overridden by explicit Chmod, so we expect 0700.
				t.Fatalf("%s perm = %o, want 700", dir, mode)
			}
		}
	}
}
