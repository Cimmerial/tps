package addonmenu

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func MagickAvailable() bool {
	return magickBin() != ""
}

func MagickCanRead(format string) bool {
	return magickFormatRW(format, true)
}

func MagickCanWrite(format string) bool {
	return magickFormatRW(format, false)
}

func magickBin() string {
	if p, err := exec.LookPath("magick"); err == nil {
		return p
	}
	if p, err := exec.LookPath("convert"); err == nil {
		return p
	}
	return ""
}

func magickFormatRW(format string, read bool) bool {
	if !MagickAvailable() {
		return false
	}
	if magickRWFromList(format, read) {
		return true
	}
	if read {
		return magickReadSmoke(format)
	}
	return magickWriteSmoke(format)
}

func magickRWFromList(format string, read bool) bool {
	names := magickFormatNames(format)
	out, err := exec.Command(magickBin(), "-list", "format").Output()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(out), "\n") {
		formatName, mode, ok := parseMagickFormatLine(line)
		if !ok {
			continue
		}
		if !formatNameMatches(formatName, names) {
			continue
		}
		if read {
			return strings.Contains(mode, "r")
		}
		return strings.Contains(mode, "w")
	}
	return false
}

func magickFormatNames(format string) []string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "jpeg":
		return []string{"JPEG", "JPG", "JPE"}
	case "tiff":
		return []string{"TIFF", "TIF"}
	case "webp":
		return []string{"WEBP"}
	default:
		return []string{strings.ToUpper(strings.TrimSpace(format))}
	}
}

func formatNameMatches(name string, names []string) bool {
	name = strings.TrimSuffix(strings.ToUpper(strings.TrimSpace(name)), "*")
	for _, n := range names {
		if name == n {
			return true
		}
	}
	return false
}

func parseMagickFormatLine(line string) (formatName, mode string, ok bool) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", "", false
	}
	if strings.EqualFold(fields[0], "Format") || strings.HasPrefix(fields[0], "-") {
		return "", "", false
	}
	formatName = fields[0]
	for i := 1; i < len(fields); i++ {
		if isModeToken(fields[i]) {
			return formatName, fields[i], true
		}
	}
	return "", "", false
}

func isModeToken(s string) bool {
	if len(s) < 2 || len(s) > 4 {
		return false
	}
	for _, c := range s {
		if c != 'r' && c != 'w' && c != '-' && c != '+' {
			return false
		}
	}
	return strings.Contains(s, "r") || strings.Contains(s, "w")
}

func magickWriteSmoke(format string) bool {
	id := strings.ToLower(strings.TrimSpace(format))
	ext := "." + id
	if id == "jpeg" {
		ext = ".jpg"
	}
	if id == "tiff" {
		ext = ".tif"
	}
	dir, err := os.MkdirTemp("", "tps-probe-")
	if err != nil {
		return false
	}
	defer os.RemoveAll(dir)
	src := filepath.Join(dir, "in.png")
	dst := filepath.Join(dir, "out"+ext)
	if _, err := exec.Command(magickBin(), "-size", "1x1", "xc:white", src).CombinedOutput(); err != nil {
		return false
	}
	if _, err := exec.Command(magickBin(), src, dst).CombinedOutput(); err != nil {
		return false
	}
	_, err = os.Stat(dst)
	return err == nil
}

func magickReadSmoke(format string) bool {
	id := strings.ToLower(strings.TrimSpace(format))
	ext := "." + id
	if id == "jpeg" {
		ext = ".jpg"
	}
	if id == "tiff" {
		ext = ".tif"
	}
	dir, err := os.MkdirTemp("", "tps-probe-")
	if err != nil {
		return false
	}
	defer os.RemoveAll(dir)
	src := filepath.Join(dir, "in"+ext)
	dst := filepath.Join(dir, "out.png")
	if _, err := exec.Command(magickBin(), "-size", "1x1", "xc:white", src).CombinedOutput(); err != nil {
		return false
	}
	if _, err := exec.Command(magickBin(), src, dst).CombinedOutput(); err != nil {
		return false
	}
	_, err = os.Stat(dst)
	return err == nil
}

func BinaryOnPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
