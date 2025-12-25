package dirs

// Paths is an app-scoped directory set.
//
// Intentional design:
// - Uses XDG Base Directory spec defaults (if env vars not set)
// - Stores secrets (e.g. TLS keys) under DataDir/PKIDir and uses 0700/0600 perms via Ensure().
// - On macOS, still follows XDG because you requested XDG-style layout.
//
// Notes:
// - launchd plist location is NOT XDG; that lives in dirs/darwin.

type Paths struct {
	AppName string

	// XDG base dirs (resolved)
	ConfigHome string
	DataHome   string
	StateHome  string
	CacheHome  string
	RuntimeDir string // may be empty on macOS

	// App-scoped dirs
	ConfigDir string // ConfigHome/AppName
	DataDir   string // DataHome/AppName
	StateDir  string // StateHome/AppName
	CacheDir  string // CacheHome/AppName

	// Common sub-dirs (convention)
	RunDir     string // StateDir/run (portable fallback)
	LogDir     string // StateDir/logs
	PKIDir     string // DataDir/pki
	ServiceDir string // StateDir/service (templates/metadata etc.)
}

// Resolve returns XDG-based Paths for the given app name.
//
// This is a facade that currently resolves XDG paths.
// If you later want to support other strategies, add options here.
func Resolve(appName string, opts ...Option) (Paths, error) {
	cfg := resolveConfig{
		appName: appName,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return resolveXDG(cfg.appName, cfg.homeDirOverride)
}

type resolveConfig struct {
	appName         string
	homeDirOverride string
}

type Option func(*resolveConfig)

// WithHomeDir overrides the home directory used for resolution.
// Intended mainly for tests or specialized environments.
func WithHomeDir(home string) Option {
	return func(c *resolveConfig) {
		c.homeDirOverride = home
	}
}
