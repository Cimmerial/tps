package config

import (
	"fmt"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type shellBucket string

const (
	shellGeneric shellBucket = "generic"
	shellWindows shellBucket = "windows"
	shellLinux   shellBucket = "linux"
	shellDarwin  shellBucket = "darwin"
)

type shellStore struct {
	generic map[int]string
	windows map[int]string
	linux   map[int]string
	darwin  map[int]string
}

var shellFieldKey = regexp.MustCompile(`^(?:(windows|linux|darwin|mac)_)?shell(\d+)?$`)

func parseShellStore(raw map[string]interface{}) shellStore {
	st := shellStore{
		generic: map[int]string{},
		windows: map[int]string{},
		linux:   map[int]string{},
		darwin:  map[int]string{},
	}
	for k, v := range raw {
		m := shellFieldKey.FindStringSubmatch(strings.ToLower(k))
		if m == nil {
			continue
		}
		val := strings.TrimSpace(fmt.Sprint(v))
		if val == "" {
			continue
		}
		idx := 1
		if m[2] != "" {
			n, err := strconv.Atoi(m[2])
			if err != nil || n < 1 {
				continue
			}
			idx = n
		}
		plat := m[1]
		switch plat {
		case "windows":
			st.windows[idx] = val
		case "linux":
			st.linux[idx] = val
		case "darwin", "mac":
			st.darwin[idx] = val
		default:
			st.generic[idx] = val
		}
	}
	return st
}

func (s shellStore) ordered(bucket shellBucket) []string {
	m := s.bucketMap(bucket)
	if len(m) == 0 {
		return nil
	}
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, m[k])
	}
	return out
}

func (s shellStore) bucketMap(bucket shellBucket) map[int]string {
	switch bucket {
	case shellWindows:
		return s.windows
	case shellLinux:
		return s.linux
	case shellDarwin:
		return s.darwin
	default:
		return s.generic
	}
}

func (s shellStore) platformBucket(goos string) shellBucket {
	switch goos {
	case "windows":
		return shellWindows
	case "linux":
		return shellLinux
	case "darwin":
		return shellDarwin
	default:
		return shellGeneric
	}
}

func (s shellStore) resolved(goos string) []string {
	plat := s.platformBucket(goos)
	if shells := s.ordered(plat); len(shells) > 0 {
		return shells
	}
	return s.ordered(shellGeneric)
}

func shellFieldName(bucket shellBucket, idx int) string {
	var prefix string
	switch bucket {
	case shellWindows:
		prefix = "windows_"
	case shellLinux:
		prefix = "linux_"
	case shellDarwin:
		prefix = "darwin_"
	}
	if idx <= 1 {
		return prefix + "shell"
	}
	return fmt.Sprintf("%sshell%d", prefix, idx)
}

func (s shellStore) marshalFields() map[string]string {
	out := map[string]string{}
	write := func(bucket shellBucket, m map[int]string) {
		if len(m) == 0 {
			return
		}
		keys := make([]int, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			out[shellFieldName(bucket, k)] = m[k]
		}
	}
	write(shellGeneric, s.generic)
	write(shellWindows, s.windows)
	write(shellLinux, s.linux)
	write(shellDarwin, s.darwin)
	return out
}

// ShellField is one YAML shell key and its command line.
type ShellField struct {
	Key   string
	Value string
}

func shellBucketFromPlat(plat string) shellBucket {
	switch plat {
	case "windows":
		return shellWindows
	case "linux":
		return shellLinux
	case "darwin", "mac":
		return shellDarwin
	default:
		return shellGeneric
	}
}

// ParseShellYAMLKey parses a shell YAML key (e.g. shell, shell2, darwin_shell).
func ParseShellYAMLKey(key string) (bucket shellBucket, idx int, ok bool) {
	m := shellFieldKey.FindStringSubmatch(strings.ToLower(strings.TrimSpace(key)))
	if m == nil {
		return shellGeneric, 0, false
	}
	idx = 1
	if m[2] != "" {
		n, err := strconv.Atoi(m[2])
		if err != nil || n < 1 {
			return shellGeneric, 0, false
		}
		idx = n
	}
	return shellBucketFromPlat(m[1]), idx, true
}

func (c *Command) ensureShellStore() *shellStore {
	if c == nil {
		return nil
	}
	if c.shells != nil {
		return c.shells
	}
	st := shellStore{
		generic: map[int]string{},
		windows: map[int]string{},
		linux:   map[int]string{},
		darwin:  map[int]string{},
	}
	if strings.TrimSpace(c.Shell) != "" {
		st.generic[1] = strings.TrimSpace(c.Shell)
	}
	c.shells = &st
	return c.shells
}

func (s *shellStore) pruneEmpty() {
	if s == nil {
		return
	}
	for _, m := range []map[int]string{s.generic, s.windows, s.linux, s.darwin} {
		for k, v := range m {
			if strings.TrimSpace(v) == "" {
				delete(m, k)
			}
		}
	}
}

func (s *shellStore) isEmpty() bool {
	if s == nil {
		return true
	}
	return len(s.generic) == 0 && len(s.windows) == 0 && len(s.linux) == 0 && len(s.darwin) == 0
}

// AllShellFields lists every shell key on the command in stable order.
func (c *Command) AllShellFields() []ShellField {
	if c == nil {
		return nil
	}
	var out []ShellField
	appendBucket := func(bucket shellBucket) {
		st := c.shells
		if st == nil {
			return
		}
		m := st.bucketMap(bucket)
		if len(m) == 0 {
			return
		}
		keys := make([]int, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			out = append(out, ShellField{
				Key:   shellFieldName(bucket, k),
				Value: m[k],
			})
		}
	}
	if c.shells != nil {
		appendBucket(shellGeneric)
		appendBucket(shellDarwin)
		appendBucket(shellLinux)
		appendBucket(shellWindows)
		return out
	}
	if strings.TrimSpace(c.Shell) != "" {
		return []ShellField{{Key: "shell", Value: c.Shell}}
	}
	return nil
}

// HasAnyShellLine reports whether the command has at least one shell line configured.
func (c *Command) HasAnyShellLine() bool {
	return len(c.AllShellFields()) > 0
}

func (c *Command) syncShellFromStore() {
	if c == nil {
		return
	}
	if c.shells == nil || c.shells.isEmpty() {
		c.Shell = ""
		c.shells = nil
		return
	}
	resolved := c.shells.resolved(runtime.GOOS)
	if len(resolved) > 0 {
		c.Shell = resolved[0]
	} else if g := c.shells.ordered(shellGeneric); len(g) > 0 {
		c.Shell = g[0]
	} else {
		c.Shell = ""
	}
}

// SetShellByKey sets one shell YAML key without removing others.
func (c *Command) SetShellByKey(key, value string) error {
	bucket, idx, ok := ParseShellYAMLKey(key)
	if !ok {
		return fmt.Errorf("invalid shell key %q", key)
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("shell value is required")
	}
	st := c.ensureShellStore()
	st.bucketMap(bucket)[idx] = value
	c.syncShellFromStore()
	return nil
}

// RemoveShellByKey removes one shell YAML key.
func (c *Command) RemoveShellByKey(key string) error {
	bucket, idx, ok := ParseShellYAMLKey(key)
	if !ok {
		return fmt.Errorf("invalid shell key %q", key)
	}
	if c.shells == nil {
		if key == "shell" && strings.TrimSpace(c.Shell) != "" {
			c.Shell = ""
			return nil
		}
		return nil
	}
	delete(c.shells.bucketMap(bucket), idx)
	c.shells.pruneEmpty()
	if c.shells.isEmpty() {
		c.shells = nil
		c.Shell = ""
	} else {
		c.syncShellFromStore()
	}
	return nil
}

// ShellPlatform selects which shell key family to edit.
type ShellPlatform int

const (
	ShellPlatformGeneric ShellPlatform = 0
	ShellPlatformDarwin  ShellPlatform = 1
	ShellPlatformLinux   ShellPlatform = 2
	ShellPlatformWindows ShellPlatform = 3
)

func shellBucketFromPlatform(p ShellPlatform) shellBucket {
	switch p {
	case ShellPlatformDarwin:
		return shellDarwin
	case ShellPlatformLinux:
		return shellLinux
	case ShellPlatformWindows:
		return shellWindows
	default:
		return shellGeneric
	}
}

// ShellFieldPlatform classifies a shell YAML key by platform.
func ShellFieldPlatform(key string) ShellPlatform {
	b, _, ok := ParseShellYAMLKey(key)
	if !ok {
		return ShellPlatformGeneric
	}
	switch b {
	case shellDarwin:
		return ShellPlatformDarwin
	case shellLinux:
		return ShellPlatformLinux
	case shellWindows:
		return ShellPlatformWindows
	default:
		return ShellPlatformGeneric
	}
}

// NextShellKey returns the YAML key for the next free index in a platform bucket.
func (c *Command) NextShellKey(platform ShellPlatform) string {
	bucket := shellBucketFromPlatform(platform)
	st := c.ensureShellStore()
	m := st.bucketMap(bucket)
	max := 0
	for k := range m {
		if k > max {
			max = k
		}
	}
	return shellFieldName(bucket, max+1)
}

// ResolvedShells returns shell lines for this OS (platform-specific keys, else shell/shell1/…).
func (c *Command) ResolvedShells() []string {
	if c == nil {
		return nil
	}
	if c.shells != nil {
		return c.shells.resolved(runtime.GOOS)
	}
	if strings.TrimSpace(c.Shell) != "" {
		return []string{c.Shell}
	}
	return nil
}
