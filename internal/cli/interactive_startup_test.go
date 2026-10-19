package cli

import (
	"os/exec"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/project"
)

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	cmds := [][]string{
		{"git", "-C", dir, "init"},
		{"git", "-C", dir, "config", "user.email", "tps@test"},
		{"git", "-C", dir, "config", "user.name", "tps"},
	}
	for _, c := range cmds {
		if out, err := exec.Command(c[0], c[1:]...).CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
}

func TestResolveInteractiveProjectNoAutoLink(t *testing.T) {
	linked := t.TempDir()
	initGitRepo(t, linked)
	linkedClean, err := project.ProjectRootFromCwd(linked)
	if err != nil {
		t.Fatal(err)
	}

	unlinked := t.TempDir()
	initGitRepo(t, unlinked)

	cfg := &config.Config{
		Projects: []config.Project{
			{ID: "existing", Name: "existing", Path: linkedClean},
		},
	}

	choice, err := resolveInteractiveProject(cfg, unlinked)
	if err != nil {
		t.Fatal(err)
	}
	if choice.NoProjects || choice.ProjectID != "" || !choice.Pick {
		t.Fatalf("want pick without auto-link: %+v", choice)
	}
	if len(cfg.Projects) != 1 {
		t.Fatalf("config mutated: %d projects", len(cfg.Projects))
	}
}

func TestResolveInteractiveProjectFromCwd(t *testing.T) {
	root := t.TempDir()
	initGitRepo(t, root)
	rootClean, err := project.ProjectRootFromCwd(root)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Projects: []config.Project{
			{ID: "p1", Name: "app", Path: rootClean},
		},
	}

	choice, err := resolveInteractiveProject(cfg, root)
	if err != nil {
		t.Fatal(err)
	}
	if choice.Pick || choice.NoProjects || choice.ProjectID != "p1" {
		t.Fatalf("want linked project: %+v", choice)
	}
}

func TestResolveInteractiveProjectNoProjects(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{}
	choice, err := resolveInteractiveProject(cfg, dir)
	if err != nil {
		t.Fatal(err)
	}
	if !choice.NoProjects || choice.Pick || choice.ProjectID != "" {
		t.Fatalf("want no projects: %+v", choice)
	}
}

func TestResolveInteractiveProjectHealByRemote(t *testing.T) {
	oldRoot := t.TempDir()
	initGitRepo(t, oldRoot)
	if out, err := exec.Command("git", "-C", oldRoot, "remote", "add", "origin", "git@example.com:org/repo.git").CombinedOutput(); err != nil {
		t.Fatalf("%v: %s", err, out)
	}

	newRoot := t.TempDir()
	initGitRepo(t, newRoot)
	if out, err := exec.Command("git", "-C", newRoot, "remote", "add", "origin", "git@example.com:org/repo.git").CombinedOutput(); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	newClean, err := project.ProjectRootFromCwd(newRoot)
	if err != nil {
		t.Fatal(err)
	}

	oldClean, err := project.ProjectRootFromCwd(oldRoot)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Projects: []config.Project{
			{ID: "p1", Name: "app", Path: oldClean, GitRemote: "git@example.com:org/repo.git"},
		},
	}

	choice, err := resolveInteractiveProject(cfg, newRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !choice.SaveConfig || choice.ProjectID != "p1" || choice.Pick {
		t.Fatalf("want healed match: %+v", choice)
	}
	healed, err := config.CleanPath(cfg.Projects[0].Path)
	if err != nil {
		t.Fatal(err)
	}
	if healed != newClean {
		t.Fatalf("path not healed: %q want %q", healed, newClean)
	}
}
