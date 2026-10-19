package config

import "strings"

// SetShellLine sets the single shell field and clears shell1 / platform-specific shells.
// Use when the user edits shell interactively.
func (c *Command) SetShellLine(shell string) {
	c.Shell = strings.TrimSpace(shell)
	c.shells = nil
}

// HasMultiShell reports whether the command uses shell1 or platform-specific shell keys.
func (c *Command) HasMultiShell() bool {
	return c != nil && c.shells != nil
}
