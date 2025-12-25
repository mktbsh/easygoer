package dirs

import (
	"os"
	"path/filepath"
	"strings"
)

func getenvOr(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

// cleanAbs:
// - expands leading "~/"
// - resolves relative paths against home
// - cleans the resulting path
func cleanAbs(p, home string) string {
	p = strings.TrimSpace(p)
	if strings.HasPrefix(p, "~/") {
		p = filepath.Join(home, p[2:])
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(home, p)
	}
	return filepath.Clean(p)
}
