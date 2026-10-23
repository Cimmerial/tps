package cli

import (
	"strings"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
	"gopkg.in/yaml.v3"
)

func TestSplitInputFieldsQuoted(t *testing.T) {
	fields := splitInputFields(`nc -n "push release" -a rel -c "bash scripts/push-release.sh [INPUT1]" -tags tps`)
	if len(fields) < 2 || fields[0] != "nc" {
		t.Fatalf("fields: %v", fields)
	}
	cmd, err := parseNewCommandFlags(fields[1:])
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Alias != "push release" {
		t.Fatalf("alias: %q", cmd.Alias)
	}
	if cmd.Abbrev != "rel" {
		t.Fatalf("abbrev: %q", cmd.Abbrev)
	}
	if cmd.Shell != "bash scripts/push-release.sh [INPUT1]" {
		t.Fatalf("shell: %q", cmd.Shell)
	}
	if len(cmd.Tags) != 1 || cmd.Tags[0] != "tps" {
		t.Fatalf("tags: %v", cmd.Tags)
	}
}

func TestParseNewCommandUnquoted(t *testing.T) {
	cmd, err := parseNewCommandFlags(splitInputFields(`nc -n push-release -a rel -c bash scripts/push-release.sh -tags tps`)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Alias != "push-release" {
		t.Fatalf("alias: %q", cmd.Alias)
	}
	if cmd.Shell != "bash scripts/push-release.sh" {
		t.Fatalf("shell: %q", cmd.Shell)
	}
}

func TestParseNewCommandShellKeepsUnknownFlags(t *testing.T) {
	cmd, err := parseNewCommandFlags(splitInputFields(`nc -n list -c ls -la`)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Shell != "ls -la" {
		t.Fatalf("shell: %q", cmd.Shell)
	}
	quoted, err := parseNewCommandFlags(splitInputFields(`nc -n list -c "ls -t"`)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if quoted.Shell != "ls -t" || quoted.Detach {
		t.Fatalf("quoted: shell %q detach %v", quoted.Shell, quoted.Detach)
	}
}

func TestParseNCShellFlag(t *testing.T) {
	cases := []struct {
		in   string
		plat string
		idx  int
		ok   bool
	}{
		{"-c", "", 0, true},
		{"-c2", "", 2, true},
		{"--shell3", "", 3, true},
		{"-mac", "darwin", 0, true},
		{"-darwin2", "darwin", 2, true},
		{"-mac_shell", "darwin", 0, true},
		{"-linux", "linux", 0, true},
		{"-linux_shell2", "linux", 2, true},
		{"-win", "windows", 0, true},
		{"-windows", "windows", 0, true},
		{"-windows_shell2", "windows", 2, true},
		{"-close", "", 0, false},
		{"-capture", "", 0, false},
		{"-H", "", 0, false},
		{"-help", "", 0, false},
		{"-fg", "", 0, false},
		{"-desc", "", 0, false},
	}
	for _, tc := range cases {
		plat, idx, ok := parseNCShellFlag(tc.in)
		if ok != tc.ok || plat != tc.plat || idx != tc.idx {
			t.Fatalf("%s: got plat %q idx %d ok %v", tc.in, plat, idx, ok)
		}
	}
}

func shellFieldMap(c config.Command) map[string]string {
	out := map[string]string{}
	for _, f := range c.AllShellFields() {
		out[f.Key] = f.Value
	}
	return out
}

func TestParseNewCommandPlatformShells(t *testing.T) {
	line := `nc -n build -a b -desc "build app" -mac ./scripts/build-mac.sh -linux ./scripts/build.sh -win "powershell -File scripts/build.ps1" -t -close`
	cmd, err := parseNewCommandFlags(splitInputFields(line)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Description != "build app" {
		t.Fatalf("desc: %q", cmd.Description)
	}
	if !cmd.Detach || !cmd.CloseOnComplete {
		t.Fatalf("detach %v close %v", cmd.Detach, cmd.CloseOnComplete)
	}
	got := shellFieldMap(cmd)
	if got["darwin_shell"] != "./scripts/build-mac.sh" || got["linux_shell"] != "./scripts/build.sh" || got["windows_shell"] != "powershell -File scripts/build.ps1" {
		t.Fatalf("shells: %#v", got)
	}

	raw, err := yaml.Marshal(cmd)
	if err != nil {
		t.Fatal(err)
	}
	var back config.Command
	if err := yaml.Unmarshal(raw, &back); err != nil {
		t.Fatal(err)
	}
	again := shellFieldMap(back)
	if again["darwin_shell"] != got["darwin_shell"] || again["linux_shell"] != got["linux_shell"] || again["windows_shell"] != got["windows_shell"] {
		t.Fatalf("round trip: %#v\n%s", again, raw)
	}
}

func TestParseNewCommandRepeatedAndNumberedShells(t *testing.T) {
	line := `nc -n stack -a st -c "npm run dev" -c "npx cloudflared" -win npm run dev:win -win2 cloudflared -t`
	cmd, err := parseNewCommandFlags(splitInputFields(line)[1:])
	if err != nil {
		t.Fatal(err)
	}
	got := shellFieldMap(cmd)
	if got["shell"] != "npm run dev" || got["shell2"] != "npx cloudflared" {
		t.Fatalf("generic: %#v", got)
	}
	if got["windows_shell"] != "npm run dev:win" || got["windows_shell2"] != "cloudflared" {
		t.Fatalf("windows: %#v", got)
	}
	if !cmd.Detach {
		t.Fatal("expected detach")
	}

	only2, err := parseNewCommandFlags(splitInputFields(`nc -n gap -shell2 echo two`)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if shellFieldMap(only2)["shell2"] != "echo two" {
		t.Fatalf("shell2: %#v", shellFieldMap(only2))
	}

	_, err = parseNewCommandFlags(splitInputFields(`nc -n twice -c2 echo -shell2 echo`)[1:])
	if err == nil || !strings.Contains(err.Error(), "set twice") {
		t.Fatalf("expected duplicate shell error, got %v", err)
	}
}

func TestParseNewCommandOptions(t *testing.T) {
	line := `nc -n tunnel -a cf -c cloudflared -m trycloudflare -actions print, clipboard, sms -tags admin, git -precursor fas -H`
	cmd, err := parseNewCommandFlags(splitInputFields(line)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if !cmd.Handoff || cmd.Detach {
		t.Fatalf("handoff %v detach %v", cmd.Handoff, cmd.Detach)
	}
	if cmd.Capture != "trycloudflare" {
		t.Fatalf("capture: %q", cmd.Capture)
	}
	if strings.Join(cmd.CaptureActions, ",") != "print,clipboard,sms" {
		t.Fatalf("actions: %#v", cmd.CaptureActions)
	}
	if len(cmd.Tags) != 2 || cmd.Tags[0] != "admin" || cmd.Tags[1] != "git" {
		t.Fatalf("tags: %#v", cmd.Tags)
	}
	if cmd.Precursor != "fas" {
		t.Fatalf("precursor: %q", cmd.Precursor)
	}

	fg, err := parseNewCommandFlags(splitInputFields(`nc -n stay -c echo -fg`)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if !fg.ForceForeground || fg.Detach {
		t.Fatalf("foreground %+v", fg)
	}

	menu, err := parseNewCommandFlags(splitInputFields(`nc -n upt -a upt -menu menus/upt.yaml`)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if menu.Menu != "menus/upt.yaml" || menu.HasAnyShellLine() {
		t.Fatalf("menu: %+v", menu)
	}

	handler, err := parseNewCommandFlags(splitInputFields(`nc -n sym -handler sym`)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if handler.Handler != "sym" {
		t.Fatalf("handler: %q", handler.Handler)
	}
}

func TestParseNewCommandRejectsConflicts(t *testing.T) {
	cases := []string{
		`nc -n x -c echo -fg -t`,
		`nc -n x -c echo -fg -m url`,
		`nc -n x -c echo -fg -close`,
		`nc -n x -c echo -fg -H`,
		`nc -n x -c echo -close -H`,
		`nc -n x -actions print`,
		`nc -n only`,
		`nc -n x -c echo -actions bogus`,
	}
	for _, line := range cases {
		_, err := parseNewCommandFlags(splitInputFields(line)[1:])
		if err == nil {
			t.Fatalf("expected error for %s", line)
		}
	}
}

func TestParseNewCommandUsage(t *testing.T) {
	_, err := parseNewCommandFlags(nil)
	if err == nil || !strings.Contains(err.Error(), "-win") || !strings.Contains(err.Error(), "-mac") {
		t.Fatalf("usage: %v", err)
	}
	_, err = parseNewCommandFlags([]string{"-help"})
	if err == nil || !strings.Contains(err.Error(), "-linux") {
		t.Fatalf("help: %v", err)
	}
}
