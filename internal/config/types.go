package config

type Config struct {
	Projects []Project `yaml:"projects"`
}

type Project struct {
	ID           string         `yaml:"id"`
	Name         string         `yaml:"name"`
	Path         string         `yaml:"path"`
	GitRemote    string         `yaml:"git_remote,omitempty"`
	Commands     []Command      `yaml:"commands,omitempty"`
	CommandUsage map[string]int `yaml:"command_usage,omitempty"`
	MetaUsage    map[string]int `yaml:"meta_usage,omitempty"`
}

type Command struct {
	Alias       string `yaml:"alias"`
	Abbrev      string `yaml:"-"`
	Shell       string `yaml:"shell,omitempty"`
	Description string `yaml:"description,omitempty"`
	Detach          bool `yaml:"detach,omitempty"`
	ForceForeground bool `yaml:"foreground,omitempty"` // skip auto-detach heuristics
	CloseOnComplete bool `yaml:"close_on_complete,omitempty"` // close detached terminal tab when the shell command exits
	Capture        string   `yaml:"capture,omitempty"`         // regex or preset (url, trycloudflare, …)
	CaptureActions []string `yaml:"capture_actions,omitempty"` // print, clipboard, sms
	Tags        []string `yaml:"tags,omitempty"`    // first tag = menu folder (e.g. admin, git)
	Precursor   string   `yaml:"precursor,omitempty"` // alias or abbrev of command to run first
	Handler     string   `yaml:"handler,omitempty"`   // in-process builtin id (no shell)
	Menu        string   `yaml:"menu,omitempty"`      // addon tree menu YAML (relative to project root)

	shells *shellStore `yaml:"-"`
}
