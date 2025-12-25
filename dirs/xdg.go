package dirs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func resolveXDG(appName string, homeOverride string) (Paths, error) {
	if appName == "" {
		return Paths{}, errors.New("appName is required")
	}

	home := homeOverride
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil || home == "" {
			return Paths{}, fmt.Errorf("failed to get home dir: %w", err)
		}
	}

	cfgHome := getenvOr("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	dataHome := getenvOr("XDG_DATA_HOME", filepath.Join(home, ".local", "share"))
	stateHome := getenvOr("XDG_STATE_HOME", filepath.Join(home, ".local", "state"))
	cacheHome := getenvOr("XDG_CACHE_HOME", filepath.Join(home, ".cache"))
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR") // may be empty on macOS

	cfgHome = cleanAbs(cfgHome, home)
	dataHome = cleanAbs(dataHome, home)
	stateHome = cleanAbs(stateHome, home)
	cacheHome = cleanAbs(cacheHome, home)
	if runtimeDir != "" {
		runtimeDir = cleanAbs(runtimeDir, home)
	}

	p := Paths{
		AppName:    appName,
		ConfigHome: cfgHome,
		DataHome:   dataHome,
		StateHome:  stateHome,
		CacheHome:  cacheHome,
		RuntimeDir: runtimeDir,
		ConfigDir:  filepath.Join(cfgHome, appName),
		DataDir:    filepath.Join(dataHome, appName),
		StateDir:   filepath.Join(stateHome, appName),
		CacheDir:   filepath.Join(cacheHome, appName),
	}

	// Portable defaults
	p.RunDir = filepath.Join(p.StateDir, "run")
	p.LogDir = filepath.Join(p.StateDir, "logs")
	p.PKIDir = filepath.Join(p.DataDir, "pki")
	p.ServiceDir = filepath.Join(p.StateDir, "service")

	return p, nil
}
