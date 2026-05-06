package tmux

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/masaway/muxflow/internal/config"
)

func TestKillSessionUsesConfiguredSocket(t *testing.T) {
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not installed")
	}

	const socketName = "muxflow-test-kill-socket"
	const sessionName = "muxflow-test-kill-session"

	SetSocket(socketName)
	defer SetSocket("")
	defer exec.Command("tmux", "-L", socketName, "kill-server").Run()

	if out, err := exec.Command("tmux", "-L", socketName, "new-session", "-d", "-s", sessionName).CombinedOutput(); err != nil {
		t.Fatalf("new-session failed: %v: %s", err, out)
	}

	if err := KillSession(sessionName); err != nil {
		t.Fatalf("KillSession failed: %v", err)
	}

	if out, err := exec.Command("tmux", "-L", socketName, "has-session", "-t", sessionName).CombinedOutput(); err == nil {
		t.Fatalf("session still exists: %s", out)
	}
}

func TestCreateSessionUsesReturnedTargetsForPaneDirectories(t *testing.T) {
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not installed")
	}

	const socketName = "muxflow-test-pane-targets"
	const sessionName = "muxflow-test-pane-dirs"

	root := t.TempDir()
	dirs := []string{"w0a", "w0b", "w1a", "w1b", "w2a", "w2b", "w3a", "w3b"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	SetSocket(socketName)
	defer SetSocket("")
	defer exec.Command("tmux", "-L", socketName, "kill-server").Run()

	project := config.Project{
		Name: sessionName,
		Path: root,
		Windows: []config.Window{
			{Name: "w0", Layout: "even-horizontal", Panes: []config.Pane{{Dir: "w0a"}, {Dir: "w0b"}}},
			{Name: "w1", Layout: "even-horizontal", Panes: []config.Pane{{Dir: "w1a"}, {Dir: "w1b"}}},
			{Name: "w2", Layout: "even-horizontal", Panes: []config.Pane{{Dir: "w2a"}, {Dir: "w2b"}}},
			{Name: "w3", Layout: "even-horizontal", Panes: []config.Pane{{Dir: "w3a"}, {Dir: "w3b"}}},
		},
	}

	created, err := CreateSession(&project, false)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	if !created {
		t.Fatal("CreateSession returned created=false")
	}

	code, out := runCmd("list-panes", "-t", sessionName, "-s", "-F", "#{window_name}|#{pane_index}|#{pane_current_path}")
	if code != 0 {
		t.Fatalf("list-panes failed")
	}

	got := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue
		}
		got[parts[0]+"|"+parts[1]] = filepath.Clean(parts[2])
	}

	want := map[string]string{
		"w0|0": filepath.Join(root, "w0a"),
		"w0|1": filepath.Join(root, "w0b"),
		"w1|0": filepath.Join(root, "w1a"),
		"w1|1": filepath.Join(root, "w1b"),
		"w2|0": filepath.Join(root, "w2a"),
		"w2|1": filepath.Join(root, "w2b"),
		"w3|0": filepath.Join(root, "w3a"),
		"w3|1": filepath.Join(root, "w3b"),
	}
	for target, expected := range want {
		if got[target] != filepath.Clean(expected) {
			t.Fatalf("%s cwd = %q, want %q\nall panes:\n%s", target, got[target], filepath.Clean(expected), out)
		}
	}

	code, optionOut := runCmd("show-window-options", "-v", "-t", sessionName+":w3", "@muxflow_window_dir")
	if code != 0 {
		t.Fatalf("@muxflow_window_dir was not set")
	}
	if filepath.Clean(optionOut) != filepath.Join(root, "w3a") {
		t.Fatalf("@muxflow_window_dir = %q, want %q", optionOut, filepath.Join(root, "w3a"))
	}

	runCmd("select-window", "-t", sessionName+":w3")
	runCmd("select-pane", "-t", sessionName+":w3.1")
	code, _ = runCmd("split-window", "-c", "#{?@muxflow_window_dir,#{@muxflow_window_dir},#{pane_current_path}}")
	if code != 0 {
		t.Fatalf("split-window with @muxflow_window_dir failed")
	}
	code, out = runCmd("display-message", "-p", "-t", sessionName+":w3.2", "#{pane_current_path}")
	if code != 0 {
		t.Fatalf("display-message for new pane failed")
	}
	if filepath.Clean(out) != filepath.Join(root, "w3a") {
		t.Fatalf("new pane cwd = %q, want %q", out, filepath.Join(root, "w3a"))
	}
}
