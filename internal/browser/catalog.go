package browser

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ID is a stable preference identifier (see All).
type ID string

const (
	IDDefault ID = "default"
	IDSafari  ID = "safari"
	IDChrome  ID = "chrome"
	IDFirefox ID = "firefox"
	IDEdge    ID = "edge"
	IDBrave   ID = "brave"
	IDArc     ID = "arc"
	IDChromium ID = "chromium"
	IDOpera   ID = "opera"
	IDVivaldi ID = "vivaldi"
)

// Profile describes a supported desktop browser.
type Profile struct {
	ID          ID
	Label       string
	BinNames    []string
	BinPaths    []string
	MacAppNames []string
	WinRelPaths []string // under Program Files / LocalAppData
}

// All lists browsers users can pick in options (default is separate).
func All() []Profile {
	home := os.Getenv("HOME")
	programFiles := os.Getenv("ProgramFiles")
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	return []Profile{
		{
			ID:          IDSafari,
			Label:       "Safari",
			MacAppNames: []string{"Safari"},
		},
		{
			ID:          IDChrome,
			Label:       "Google Chrome",
			BinNames:    []string{"google-chrome", "google-chrome-stable", "chrome"},
			BinPaths: []string{
				"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
				filepath.Join(home, "Applications", "Google Chrome.app", "Contents", "MacOS", "Google Chrome"),
			},
			WinRelPaths: winPaths(programFiles, programFilesX86,
				"Google/Chrome/Application/chrome.exe",
				"Google/Chrome/Application/chrome.exe",
			),
			MacAppNames: []string{"Google Chrome"},
		},
		{
			ID:          IDFirefox,
			Label:       "Firefox",
			BinNames:    []string{"firefox"},
			BinPaths: []string{
				"/Applications/Firefox.app/Contents/MacOS/firefox",
				filepath.Join(home, "Applications", "Firefox.app", "Contents", "MacOS", "firefox"),
			},
			WinRelPaths: winPaths(programFiles, programFilesX86,
				"Mozilla Firefox/firefox.exe",
				"Mozilla Firefox/firefox.exe",
			),
			MacAppNames: []string{"Firefox"},
		},
		{
			ID:          IDEdge,
			Label:       "Microsoft Edge",
			BinNames:    []string{"microsoft-edge", "msedge"},
			BinPaths: []string{
				"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
				filepath.Join(home, "Applications", "Microsoft Edge.app", "Contents", "MacOS", "Microsoft Edge"),
			},
			WinRelPaths: winPaths(programFiles, programFilesX86,
				"Microsoft/Edge/Application/msedge.exe",
				"Microsoft/Edge/Application/msedge.exe",
			),
			MacAppNames: []string{"Microsoft Edge"},
		},
		{
			ID:          IDBrave,
			Label:       "Brave",
			BinNames:    []string{"brave-browser", "brave"},
			BinPaths: []string{
				"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
				filepath.Join(home, "Applications", "Brave Browser.app", "Contents", "MacOS", "Brave Browser"),
			},
			WinRelPaths: winPaths(programFiles, programFilesX86,
				"BraveSoftware/Brave-Browser/Application/brave.exe",
				"BraveSoftware/Brave-Browser/Application/brave.exe",
			),
			MacAppNames: []string{"Brave Browser"},
		},
		{
			ID:          IDArc,
			Label:       "Arc",
			BinPaths: []string{
				"/Applications/Arc.app/Contents/MacOS/Arc",
				filepath.Join(home, "Applications", "Arc.app", "Contents", "MacOS", "Arc"),
			},
			MacAppNames: []string{"Arc"},
		},
		{
			ID:          IDChromium,
			Label:       "Chromium",
			BinNames:    []string{"chromium", "chromium-browser"},
			BinPaths: []string{
				"/Applications/Chromium.app/Contents/MacOS/Chromium",
				filepath.Join(home, "Applications", "Chromium.app", "Contents", "MacOS", "Chromium"),
			},
			MacAppNames: []string{"Chromium"},
		},
		{
			ID:          IDOpera,
			Label:       "Opera",
			BinNames:    []string{"opera"},
			BinPaths: []string{
				"/Applications/Opera.app/Contents/MacOS/Opera",
				filepath.Join(home, "Applications", "Opera.app", "Contents", "MacOS", "Opera"),
			},
			WinRelPaths: winPaths(programFiles, programFilesX86,
				"Opera/launcher.exe",
				"Opera/launcher.exe",
			),
			MacAppNames: []string{"Opera"},
		},
		{
			ID:          IDVivaldi,
			Label:       "Vivaldi",
			BinNames:    []string{"vivaldi"},
			BinPaths: []string{
				"/Applications/Vivaldi.app/Contents/MacOS/Vivaldi",
				filepath.Join(home, "Applications", "Vivaldi.app", "Contents", "MacOS", "Vivaldi"),
			},
			WinRelPaths: winPaths(programFiles, programFilesX86,
				"Vivaldi/Application/vivaldi.exe",
				"Vivaldi/Application/vivaldi.exe",
			),
			MacAppNames: []string{"Vivaldi"},
		},
	}
}

func winPaths(pf, pfx86, pfRel, pfx86Rel string) []string {
	var out []string
	if pf != "" && pfRel != "" {
		out = append(out, filepath.Join(pf, filepath.FromSlash(pfRel)))
	}
	if pfx86 != "" && pfx86Rel != "" {
		out = append(out, filepath.Join(pfx86, filepath.FromSlash(pfx86Rel)))
	}
	return out
}

func ProfileByID(id ID) (Profile, bool) {
	for _, p := range All() {
		if p.ID == id {
			return p, true
		}
	}
	return Profile{}, false
}

// NormalizeID canonicalizes a preference value (empty → default).
func NormalizeID(raw string) ID {
	raw = strings.ToLower(strings.TrimSpace(raw))
	switch raw {
	case "", "default", "auto", "system":
		return IDDefault
	case "safari":
		return IDSafari
	case "chrome", "google chrome", "google-chrome":
		return IDChrome
	case "firefox":
		return IDFirefox
	case "edge", "msedge", "microsoft edge":
		return IDEdge
	case "brave", "brave browser":
		return IDBrave
	case "arc":
		return IDArc
	case "chromium":
		return IDChromium
	case "opera":
		return IDOpera
	case "vivaldi":
		return IDVivaldi
	default:
		return ID(raw)
	}
}

func (p Profile) LabelOrID() string {
	if p.Label != "" {
		return p.Label
	}
	return string(p.ID)
}

func (p Profile) ResolveBin() (bin string, ok bool) {
	for _, name := range p.BinNames {
		if looked, err := exec.LookPath(name); err == nil {
			return looked, true
		}
	}
	for _, path := range p.BinPaths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}
	for _, path := range p.WinRelPaths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}
	return "", false
}

func macAppBundlePath(name string) string {
	for _, root := range []string{"/Applications", filepath.Join(os.Getenv("HOME"), "Applications")} {
		p := filepath.Join(root, name+".app")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// ResolveMacApp returns the macOS app name for `open -a`, if installed.
func (p Profile) ResolveMacApp() (name string, ok bool) {
	for _, n := range p.MacAppNames {
		if macAppBundlePath(n) != "" {
			return n, true
		}
	}
	return "", false
}
