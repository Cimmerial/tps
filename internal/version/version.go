package version

// Version is set at link time via -ldflags on release builds; default "dev" for local go install.
// User-facing text should use Label() (module semver from the binary, not git in cwd).
var Version = "dev"

// DevVersion is the semver shown for (devel) builds. Keep in sync with release tags.
const DevVersion = "0.1.10"
