package addonmenu

import (
	"os/exec"
	"strings"
)

func MagickAvailable() bool {
	_, err := exec.LookPath("magick")
	return err == nil
}

// MagickCanRead reports whether ImageMagick lists a decoder for the format id (jpeg, heic, …).
func MagickCanRead(format string) bool {
	return magickFormatRW(format, true)
}

// MagickCanWrite reports whether ImageMagick lists an encoder for the format id.
func MagickCanWrite(format string) bool {
	return magickFormatRW(format, false)
}

func magickFormatRW(format string, read bool) bool {
	if !MagickAvailable() {
		return false
	}
	format = strings.ToUpper(strings.TrimSpace(format))
	if format == "" {
		return false
	}
	out, err := exec.Command("magick", "-list", "format").Output()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := strings.TrimSuffix(strings.ToUpper(fields[0]), "*")
		if name != format && !strings.HasPrefix(name, format) {
			continue
		}
		flags := fields[1]
		if read {
			if strings.Contains(flags, "r") {
				return true
			}
		} else if strings.Contains(flags, "w") {
			return true
		}
	}
	return false
}

func BinaryOnPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
