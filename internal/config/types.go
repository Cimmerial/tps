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
}

type Command struct {
	Alias       string `yaml:"alias"`
	Abbrev      string `yaml:"-"`
	Shell       string `yaml:"shell"`
	Description string `yaml:"description,omitempty"`
}
