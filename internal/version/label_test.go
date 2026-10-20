package version

import "testing"

func TestLabelRelease(t *testing.T) {
	old := Version
	Version = "v0.1.0"
	defer func() { Version = old }()
	if Label() != "v0.1.0" {
		t.Fatalf("Label() = %q", Label())
	}
}

func TestLabelDevUsesGitDescribe(t *testing.T) {
	oldV := Version
	oldD := describeGit
	Version = "dev"
	describeGit = func() string { return "v0.1.0-2-gabc" }
	defer func() {
		Version = oldV
		describeGit = oldD
	}()
	if Label() != "v0.1.0-2-gabc" {
		t.Fatalf("Label() = %q", Label())
	}
}
