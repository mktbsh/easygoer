package dirs

import (
	"errors"
	"fmt"
	"os"
)

// Ensure creates app directories with safe permissions.
//
// - Directories are created with 0700 (explicit chmod after MkdirAll).
// - This is a safe baseline (works for secrets too).
//
// Note: chmod may fail on some FS; this function treats it as error.
func (p Paths) Ensure() error {
	if p.AppName == "" {
		return errors.New("Ensure: Paths is not resolved (AppName empty)")
	}

	for _, dir := range []string{
		p.ConfigDir,
		p.DataDir,
		p.StateDir,
		p.CacheDir,
		p.RunDir,
		p.LogDir,
		p.PKIDir,
		p.ServiceDir,
	} {
		if err := ensureDir(dir); err != nil {
			return err
		}
	}
	return nil
}

func ensureDir(path string) error {
	if path == "" {
		return errors.New("ensureDir: empty path")
	}
	if err := os.MkdirAll(path, 0o700); err != nil {
		return fmt.Errorf("ensureDir %s: %w", path, err)
	}
	if err := os.Chmod(path, 0o700); err != nil {
		return fmt.Errorf("chmod %s: %w", path, err)
	}
	return nil
}
