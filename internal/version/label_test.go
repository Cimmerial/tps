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
	oldM := moduleVersionLabel
	Version = "dev"
	moduleVersionLabel = func() string { return "" }
	describeGit = func() string { return "v0.1.0-2-gabc" }
	defer func() {
		Version = oldV
		describeGit = oldD
		moduleVersionLabel = oldM
	}()
	if Label() != "v0.1.0-2-gabc" {
		t.Fatalf("Label() = %q", Label())
	}
}

func TestLabelIgnoresBareCommitFromGitDescribe(t *testing.T) {
	oldV := Version
	oldD := describeGit
	oldM := moduleVersionLabel
	Version = "dev"
	moduleVersionLabel = func() string { return "" }
	describeGit = func() string { return "b1eb6f0" }
	defer func() {
		Version = oldV
		describeGit = oldD
		moduleVersionLabel = oldM
	}()
	if Label() != "v0.1.0" {
		t.Fatalf("Label() = %q", Label())
	}
}
