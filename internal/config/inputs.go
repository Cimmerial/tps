package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var shellInputPlaceholder = regexp.MustCompile(`\[INPUT(\d+)(\?)?\]`)

func shellInputBounds(c *Command) (requiredMax, acceptedMax int) {
	if c == nil {
		return 0, 0
	}
	for _, line := range c.ResolvedShells() {
		for _, m := range shellInputPlaceholder.FindAllStringSubmatch(line, -1) {
			if len(m) < 2 {
				continue
			}
			n, err := strconv.Atoi(m[1])
			if err != nil || n < 1 {
				continue
			}
			optional := len(m) >= 3 && m[2] == "?"
			if n > acceptedMax {
				acceptedMax = n
			}
			if !optional && n > requiredMax {
				requiredMax = n
			}
		}
	}
	return requiredMax, acceptedMax
}

// MaxShellInputIndex is the highest required [INPUTn] index (optional [INPUTn?] does not count).
func MaxShellInputIndex(c *Command) int {
	required, _ := shellInputBounds(c)
	return required
}

// ParseShellInputTokens parses -value tokens after a command name on the interactive line.
func ParseShellInputTokens(fields []string) ([]string, error) {
	if len(fields) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if !strings.HasPrefix(f, "-") {
			return nil, fmt.Errorf("command inputs must be -value tokens (got %q)", f)
		}
		v := strings.TrimPrefix(f, "-")
		if v == "" {
			return nil, fmt.Errorf("empty command input value")
		}
		out = append(out, v)
	}
	return out, nil
}

// ValidateCommandInputs checks input count against [INPUTn] placeholders in the command shell(s).
func ValidateCommandInputs(c *Command, values []string) error {
	if c == nil {
		return fmt.Errorf("no command")
	}
	required, accepted := shellInputBounds(c)
	if accepted == 0 {
		if len(values) > 0 {
			return fmt.Errorf("command %q does not accept inputs", c.Alias)
		}
		return nil
	}
	if len(values) < required {
		return fmt.Errorf("command %q needs %d input(s) (use -value after the abbrev)", c.Alias, required)
	}
	if len(values) > accepted {
		return fmt.Errorf("command %q accepts at most %d input(s)", c.Alias, accepted)
	}
	return nil
}

func shellInputMap(values []string) map[int]string {
	m := make(map[int]string, len(values))
	for i, v := range values {
		m[i+1] = v
	}
	return m
}

// SubstituteShellInputs replaces [INPUTn] in a shell line. Returns an error if any placeholder remains.
func SubstituteShellInputs(shell string, values []string) (string, error) {
	if !shellInputPlaceholder.MatchString(shell) {
		return shell, nil
	}
	m := shellInputMap(values)
	out := shellInputPlaceholder.ReplaceAllStringFunc(shell, func(ph string) string {
		sub := shellInputPlaceholder.FindStringSubmatch(ph)
		if len(sub) < 2 {
			return ph
		}
		n, err := strconv.Atoi(sub[1])
		if err != nil || n < 1 {
			return ph
		}
		optional := len(sub) >= 3 && sub[2] == "?"
		if v, ok := m[n]; ok {
			return v
		}
		if optional {
			return ""
		}
		return ph
	})
	if shellInputPlaceholder.MatchString(out) {
		return "", fmt.Errorf("unset shell placeholder in %q", shell)
	}
	return out, nil
}

// ShellsWithInputs returns resolved shell lines with [INPUTn] replaced by values.
func (c *Command) ShellsWithInputs(values []string) ([]string, error) {
	if err := ValidateCommandInputs(c, values); err != nil {
		return nil, err
	}
	shells := c.ResolvedShells()
	if len(shells) == 0 {
		return nil, fmt.Errorf("command %q has no shell", c.Alias)
	}
	out := make([]string, len(shells))
	for i, line := range shells {
		sub, err := SubstituteShellInputs(line, values)
		if err != nil {
			return nil, err
		}
		out[i] = sub
	}
	return out, nil
}
