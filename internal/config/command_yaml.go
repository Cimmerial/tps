package config

// UnmarshalYAML accepts abbreviation (preferred) and abbrev (legacy).
func (c *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type raw struct {
		Alias        string `yaml:"alias"`
		Shell        string `yaml:"shell"`
		Description  string `yaml:"description,omitempty"`
		Abbreviation string `yaml:"abbreviation,omitempty"`
		Abbrev       string `yaml:"abbrev,omitempty"`
		Detach       bool   `yaml:"detach,omitempty"`
		Capture      string `yaml:"capture,omitempty"`
	}
	var r raw
	if err := unmarshal(&r); err != nil {
		return err
	}
	c.Alias = r.Alias
	c.Shell = r.Shell
	c.Description = r.Description
	c.Abbrev = r.Abbreviation
	if c.Abbrev == "" {
		c.Abbrev = r.Abbrev
	}
	c.Detach = r.Detach
	c.Capture = r.Capture
	return nil
}

// MarshalYAML writes abbreviation when abbrev is set.
func (c Command) MarshalYAML() (interface{}, error) {
	type out struct {
		Alias        string `yaml:"alias"`
		Shell        string `yaml:"shell"`
		Description  string `yaml:"description,omitempty"`
		Abbreviation string `yaml:"abbreviation,omitempty"`
		Detach       bool   `yaml:"detach,omitempty"`
		Capture      string `yaml:"capture,omitempty"`
	}
	row := out{
		Alias:       c.Alias,
		Shell:       c.Shell,
		Description: c.Description,
		Detach:      c.Detach,
		Capture:     c.Capture,
	}
	if c.Abbrev != "" {
		row.Abbreviation = c.Abbrev
	}
	return row, nil
}
