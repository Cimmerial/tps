package addonmenu

type MenuFile struct {
	Title string     `yaml:"title"`
	Items []MenuItem `yaml:"items"`
}

type MenuItem struct {
	ID       string      `yaml:"id"`
	Label    string      `yaml:"label"`
	Children []MenuItem  `yaml:"children"`
	Back     bool        `yaml:"back"`
	Action   *MenuAction `yaml:"action"`
	DimWhen  *DimWhen    `yaml:"dim_when"`
}

type DimWhen struct {
	Reason         string `yaml:"reason"`
	RequireBinary  string `yaml:"require_binary"`
	MagickRead     string `yaml:"magick_read"`
	MagickWrite    string `yaml:"magick_write"`
}

type MenuAction struct {
	Picks []PickSpec `yaml:"picks"`
	Shell string     `yaml:"shell"`
}

type PickSpec struct {
	Kind     string `yaml:"kind"` // file, files, folder
	Prompt   string `yaml:"prompt"`
	Optional bool   `yaml:"optional"`
}
