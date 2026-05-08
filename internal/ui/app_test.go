package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/masaway/muxflow/internal/config"
)

func TestLatestProjectConfigPrefersSavedProjectWithSameNameAndPath(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	config.SetSocket("ui-latest-project")
	defer config.SetSocket("")

	root := t.TempDir()
	stale := config.Project{
		Name: "crosslog",
		Path: root,
		Windows: []config.Window{
			{Name: "boot", Layout: "tiled", Panes: []config.Pane{{Dir: "."}}},
		},
	}
	latest := config.Project{
		Name: "crosslog",
		Path: root,
		Windows: []config.Window{
			{Name: "boot", Layout: "tiled", Panes: []config.Pane{{Dir: "."}}},
			{Name: "crosslog-back", Layout: "even-horizontal", Panes: []config.Pane{{Dir: "crosslog-back"}}},
		},
	}

	if err := config.Save(&config.Config{Projects: []config.Project{latest}}); err != nil {
		t.Fatal(err)
	}

	got := latestProjectConfig(stale)
	if len(got.Windows) != len(latest.Windows) {
		t.Fatalf("windows = %d, want %d", len(got.Windows), len(latest.Windows))
	}
	if got.Windows[1].Name != "crosslog-back" {
		t.Fatalf("window[1].Name = %q, want crosslog-back", got.Windows[1].Name)
	}
}

func TestLatestProjectConfigFallsBackToNameWhenPathChanged(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	config.SetSocket("ui-latest-project-path")
	defer config.SetSocket("")

	oldRoot := t.TempDir()
	newRoot := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(newRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := config.Save(&config.Config{Projects: []config.Project{{
		Name: "crosslog",
		Path: newRoot,
		Windows: []config.Window{
			{Name: "latest", Layout: "even-horizontal", Panes: []config.Pane{{Dir: "."}}},
		},
	}}}); err != nil {
		t.Fatal(err)
	}

	got := latestProjectConfig(config.Project{Name: "crosslog", Path: oldRoot})
	if got.Path != newRoot {
		t.Fatalf("Path = %q, want %q", got.Path, newRoot)
	}
}
