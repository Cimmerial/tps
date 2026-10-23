package config

import (
	"strings"
	"testing"
)

func TestExportDefaultFilename(t *testing.T) {
	got := ExportDefaultFilename(ExportKindProject, "My App", "alice")
	want := "tps-My-App-alice.tps.yaml"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if !strings.HasSuffix(ExportDefaultFilename(ExportKindOptions, "", "bob"), "tps-options-bob.tps.yaml") {
		t.Fatal("options filename")
	}
	if !strings.HasSuffix(ExportDefaultFilename(ExportKindEverything, "", "bob"), "tps-all-bob.tps.yaml") {
		t.Fatal("all filename")
	}
}

func TestExportV2RoundTrip(t *testing.T) {
	prefs := &Preferences{Username: "u1", Accent: "blue"}
	data, err := ExportPreferencesBundle(prefs, true)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ImportFromFile(writeTempExport(t, data))
	if err != nil {
		t.Fatal(err)
	}
	if b.Version != exportVersion2 || !b.HasPreferences() {
		t.Fatalf("bundle %+v", b)
	}
	if b.Preferences.Username != "u1" || b.Preferences.SessionProjectID != "" {
		t.Fatalf("prefs %+v", b.Preferences)
	}

	proj := Project{Name: "demo", Commands: []Command{{Alias: "run", Abbrev: "r", Shell: "npm test"}}}
	data, err = ExportProjectCommands(&proj)
	if err != nil {
		t.Fatal(err)
	}
	b, err = ImportFromFile(writeTempExport(t, data))
	if err != nil {
		t.Fatal(err)
	}
	cmds, ok := b.CommandsForImport("demo")
	if !ok || len(cmds) != 1 || cmds[0].Shell != "npm test" {
		t.Fatalf("cmds %+v ok=%v", cmds, ok)
	}
}

func TestMergeCommandsByAlias(t *testing.T) {
	p := Project{
		Name: "x",
		Commands: []Command{
			{Alias: "run", Shell: "old"},
			{Alias: "local", Shell: "stay"},
		},
	}
	imported := []Command{
		{Alias: "run", Shell: "new"},
		{Alias: "added", Shell: "y"},
	}
	MergeCommandsByAlias(&p, imported)
	if p.FindCommand("run").Shell != "new" {
		t.Fatal("expected update")
	}
	if p.FindCommand("local").Shell != "stay" {
		t.Fatal("expected keep local")
	}
	if p.FindCommand("added") == nil {
		t.Fatal("expected append")
	}
}

func TestCommandsForImportV1(t *testing.T) {
	b := ExportBundle{
		Version: 1,
		Project: &Project{
			ID:       "id1",
			Name:     "legacy",
			Commands: []Command{{Alias: "a", Shell: "s"}},
		},
	}
	cmds, ok := b.CommandsForImport("other")
	if !ok || len(cmds) != 1 {
		t.Fatalf("cmds=%v ok=%v", cmds, ok)
	}
}

func TestApplyPortablePreferencesPreservesSession(t *testing.T) {
	dst := &Preferences{SessionProjectID: "keep", Accent: "cyan"}
	src := &Preferences{SessionProjectID: "drop", Accent: "magenta", Username: "z"}
	ApplyPortablePreferences(dst, src, true)
	if dst.SessionProjectID != "keep" {
		t.Fatalf("session=%q", dst.SessionProjectID)
	}
	if dst.Accent != "magenta" || dst.Username != "z" {
		t.Fatalf("dst %+v", dst)
	}
}

func TestExportPreferencesOmitsUsername(t *testing.T) {
	prefs := &Preferences{Username: "secret", Accent: "blue"}
	data, err := ExportPreferencesBundle(prefs, false)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "secret") {
		t.Fatalf("username leaked: %s", data)
	}
	b, err := ImportFromFile(writeTempExport(t, data))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(b.Preferences.Username) != "" {
		t.Fatalf("username in bundle: %q", b.Preferences.Username)
	}
}

func TestApplyPortablePreferencesSkipsUsername(t *testing.T) {
	dst := &Preferences{Username: "keep"}
	src := &Preferences{Username: "replace", Accent: "magenta"}
	ApplyPortablePreferences(dst, src, false)
	if dst.Username != "keep" {
		t.Fatalf("username=%q", dst.Username)
	}
	if dst.Accent != "magenta" {
		t.Fatal("other fields should apply")
	}
}

func writeTempExport(t *testing.T, data []byte) string {
	path := t.TempDir() + "/export.tps.yaml"
	if err := WriteExportFile(path, data); err != nil {
		t.Fatal(err)
	}
	return path
}
