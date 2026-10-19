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
