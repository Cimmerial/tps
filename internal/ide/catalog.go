package ide

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ID is a stable preference / CLI identifier (see All).
type ID string

const (
	IDAuto          ID = "auto"
	IDCursor        ID = "cursor"
	IDVSCode        ID = "vscode"
	IDWindsurf      ID = "windsurf"
	IDVSCodium      ID = "codium"
	IDZed           ID = "zed"
	IDIntelliJ      ID = "intellij"
	IDWebStorm      ID = "webstorm"
	IDPyCharm       ID = "pycharm"
	IDGoLand        ID = "goland"
	IDRider         ID = "rider"
	IDRustRover     ID = "rustrover"
	IDCLion         ID = "clion"
	IDPhpStorm      ID = "phpstorm"
	IDFleet         ID = "fleet"
	IDAndroidStudio ID = "android_studio"
	IDSublime       ID = "sublime"
	IDNova          ID = "nova"
	IDXcode         ID = "xcode"
	IDLapce         ID = "lapce"
	IDBBEdit        ID = "bbedit"
)

// Profile describes a supported desktop editor.
type Profile struct {
	ID           ID
	Label        string
	ProcessName  string // macOS System Events application process name
	DataDirName  string // ~/Library/Application Support/<name> (VS Code family)
	BinNames     []string
	BinPaths     []string
	MacAppNames  []string // macOS open -a fallback when CLI missing
}

// All lists editors users can pick in options (auto is separate).
func All() []Profile {
	out := []Profile{
		{
			ID:          IDCursor,
			Label:       "Cursor",
			ProcessName: "Cursor",
			DataDirName: "Cursor",
			BinNames:    []string{"cursor"},
			BinPaths: []string{
				filepath.Join(os.Getenv("HOME"), "bin", "cursor"),
				"/Applications/Cursor.app/Contents/Resources/app/bin/cursor",
			},
		},
		{
			ID:          IDVSCode,
			Label:       "Visual Studio Code",
			ProcessName: "Code",
			DataDirName: "Code",
			BinNames:    []string{"code"},
			BinPaths: []string{
				"/usr/local/bin/code",
				filepath.Join(os.Getenv("HOME"), "bin", "code"),
				"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
			},
		},
		{
			ID:          IDWindsurf,
			Label:       "Windsurf",
			ProcessName: "Windsurf",
			DataDirName: "Windsurf",
			BinNames:    []string{"windsurf"},
			BinPaths: []string{
				filepath.Join(os.Getenv("HOME"), "bin", "windsurf"),
				"/Applications/Windsurf.app/Contents/Resources/app/bin/windsurf",
			},
		},
		{
			ID:          IDVSCodium,
			Label:       "VSCodium",
			ProcessName: "VSCodium",
			DataDirName: "VSCodium",
			BinNames:    []string{"codium"},
			BinPaths: []string{
				filepath.Join(os.Getenv("HOME"), "bin", "codium"),
				"/Applications/VSCodium.app/Contents/Resources/app/bin/codium",
			},
		},
		{
			ID:          IDZed,
			Label:       "Zed",
			ProcessName: "Zed",
			DataDirName: "Zed",
			BinNames:    []string{"zed", "zedit"},
			BinPaths: []string{
				filepath.Join(os.Getenv("HOME"), "bin", "zed"),
				"/Applications/Zed.app/Contents/MacOS/zed",
			},
		},
	}
	out = append(out, jetbrainsProfiles()...)
	out = append(out, []Profile{
		{
			ID:          IDSublime,
			Label:       "Sublime Text",
			ProcessName: "sublime_text",
			BinNames:    []string{"subl", "sublime"},
			BinPaths: []string{
				"/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl",
				filepath.Join(os.Getenv("HOME"), "bin", "subl"),
			},
			MacAppNames: []string{"Sublime Text"},
		},
		{
			ID:          IDNova,
			Label:       "Nova",
			ProcessName: "Nova",
			BinNames:    []string{"nova"},
			BinPaths: []string{
				"/Applications/Nova.app/Contents/MacOS/Nova",
				filepath.Join(os.Getenv("HOME"), "bin", "nova"),
			},
			MacAppNames: []string{"Nova"},
		},
		{
			ID:          IDXcode,
			Label:       "Xcode",
			ProcessName: "Xcode",
			BinNames:    []string{"xed"},
			BinPaths:    []string{"/usr/bin/xed"},
			MacAppNames: []string{"Xcode"},
		},
		{
			ID:          IDLapce,
			Label:       "Lapce",
			ProcessName: "lapce",
			BinNames:    []string{"lapce"},
			BinPaths: []string{
				filepath.Join(os.Getenv("HOME"), "bin", "lapce"),
				"/Applications/Lapce.app/Contents/MacOS/lapce",
			},
			MacAppNames: []string{"Lapce"},
		},
		{
			ID:          IDBBEdit,
			Label:       "BBEdit",
			ProcessName: "BBEdit",
			BinNames:    []string{"bbedit"},
			BinPaths: []string{
				"/Applications/BBEdit.app/Contents/Helpers/bbedit",
				"/usr/local/bin/bbedit",
			},
			MacAppNames: []string{"BBEdit"},
		},
	}...)
	return out
}

func jetbrainsProfiles() []Profile {
	return []Profile{
		jetbrainsProfile(IDIntelliJ, "IntelliJ IDEA", "idea", "idea",
			"IntelliJ IDEA", "IntelliJ IDEA CE", "IntelliJ IDEA Community Edition"),
		jetbrainsProfile(IDWebStorm, "WebStorm", "webstorm", "webstorm", "WebStorm"),
		jetbrainsProfile(IDPyCharm, "PyCharm", "pycharm", "pycharm",
			"PyCharm", "PyCharm CE", "PyCharm Community Edition"),
		jetbrainsProfile(IDGoLand, "GoLand", "goland", "goland", "GoLand"),
		jetbrainsProfile(IDRider, "Rider", "rider", "rider", "Rider"),
		jetbrainsProfile(IDRustRover, "RustRover", "RustRover", "rustrover", "RustRover"),
		jetbrainsProfile(IDCLion, "CLion", "clion", "clion", "CLion"),
		jetbrainsProfile(IDPhpStorm, "PhpStorm", "phpstorm", "phpstorm", "PhpStorm"),
		jetbrainsProfile(IDFleet, "Fleet", "fleet", "fleet", "Fleet"),
		jetbrainsProfile(IDAndroidStudio, "Android Studio", "studio", "studio", "Android Studio"),
	}
}

func jetbrainsProfile(id ID, label, processName, cli string, macApps ...string) Profile {
	home := os.Getenv("HOME")
	toolbox := filepath.Join(home, "Library", "Application Support", "JetBrains", "Toolbox", "scripts", cli)
	paths := []string{toolbox, filepath.Join(home, "bin", cli)}
	for _, app := range macApps {
		paths = append(paths,
			filepath.Join("/Applications", app+".app", "Contents", "MacOS", cli),
			filepath.Join(home, "Applications", app+".app", "Contents", "MacOS", cli),
		)
	}
	return Profile{
		ID:          id,
		Label:       label,
		ProcessName: processName,
		BinNames:    []string{cli},
		BinPaths:    paths,
		MacAppNames: macApps,
	}
}

func ProfileByID(id ID) (Profile, bool) {
	for _, p := range All() {
		if p.ID == id {
			return p, true
		}
	}
	return Profile{}, false
}

// NormalizeID canonicalizes a preference value (empty → auto).
func NormalizeID(raw string) ID {
	raw = strings.ToLower(strings.TrimSpace(raw))
	switch raw {
	case "", "auto", "detect", "default":
		return IDAuto
	case "cursor":
		return IDCursor
	case "vscode", "code", "vs code", "visual studio code":
		return IDVSCode
	case "windsurf":
		return IDWindsurf
	case "codium", "vscodium":
		return IDVSCodium
	case "zed", "zedit":
		return IDZed
	case "intellij", "idea", "intellij idea":
		return IDIntelliJ
	case "webstorm":
		return IDWebStorm
	case "pycharm":
		return IDPyCharm
	case "goland", "go land":
		return IDGoLand
	case "rider":
		return IDRider
	case "rustrover", "rust rover":
		return IDRustRover
	case "clion":
		return IDCLion
	case "phpstorm", "php storm":
		return IDPhpStorm
	case "fleet":
		return IDFleet
	case "android_studio", "android studio", "androidstudio":
		return IDAndroidStudio
	case "sublime", "sublime text", "subl":
		return IDSublime
	case "nova":
		return IDNova
	case "xcode", "xed":
		return IDXcode
	case "lapce":
		return IDLapce
	case "bbedit", "bb edit":
		return IDBBEdit
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

// ResolveBin returns an installed launcher for this profile, if any.
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

// ResolveLauncher returns a CLI binary or macOS app name to open the project.
func (p Profile) ResolveLauncher() (bin string, macApp string, ok bool) {
	if b, ok := p.ResolveBin(); ok {
		return b, "", true
	}
	for _, name := range p.MacAppNames {
		if macAppBundlePath(name) != "" {
			return "", name, true
		}
	}
	return "", "", false
}

// DetectInstalled picks the first installed profile (VS Code before Cursor, matching config editor detection).
func DetectInstalled() (Profile, bool) {
	order := []ID{
		IDVSCode, IDCursor, IDWindsurf, IDVSCodium, IDZed,
		IDIntelliJ, IDWebStorm, IDPyCharm, IDGoLand, IDAndroidStudio,
		IDSublime, IDNova, IDXcode,
	}
	for _, id := range order {
		p, ok := ProfileByID(id)
		if !ok {
			continue
		}
		if _, _, ok := p.ResolveLauncher(); ok {
			return p, true
		}
	}
	return Profile{}, false
}

func profileFromEditorBin(bin string) (Profile, bool) {
	base := strings.ToLower(filepath.Base(bin))
	for _, p := range All() {
		for _, name := range p.BinNames {
			if base == name {
				return p, true
			}
		}
	}
	return Profile{}, false
}
