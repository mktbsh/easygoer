package dirs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func resolveXDG(appName string, baseDir string) (Paths, error) {
	if appName == "" {
		return Paths{}, errors.New("appName is required")
	}

	home := baseDir
	useEnv := baseDir == ""

	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil || home == "" {
			return Paths{}, fmt.Errorf("failed to get home dir: %w", err)
		}
	}

	home = filepath.Clean(home)

	cfgHome := filepath.Join(home, ".config")
	dataHome := filepath.Join(home, ".local", "share")
	stateHome := filepath.Join(home, ".local", "state")
	cacheHome := filepath.Join(home, ".cache")
	runtimeDir := ""

	if useEnv {
		cfgHome = getenvOr("XDG_CONFIG_HOME", cfgHome)
		dataHome = getenvOr("XDG_DATA_HOME", dataHome)
		stateHome = getenvOr("XDG_STATE_HOME", stateHome)
		cacheHome = getenvOr("XDG_CACHE_HOME", cacheHome)
		runtimeDir = os.Getenv("XDG_RUNTIME_DIR")
	}

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
