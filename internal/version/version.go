package version

// Version is set at link time via -ldflags on release builds; default "dev" for local go install.
// User-facing text should use Label() so dev clones can show git describe.
var Version = "dev"
