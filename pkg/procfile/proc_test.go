package procfile

import (
	"strings"
	"testing"
)

func TestNewCommand(t *testing.T) {
	proc := New("test", []string{"echo", "hello"})
	cmd := proc.NewCommand()

	if !strings.HasSuffix(cmd.Path, "echo") {
		t.Errorf("Expected command path to end with 'echo', got '%s'", cmd.Path)
	}

	expectedArgs := []string{"echo", "hello"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(cmd.Args))
	}

	for i, expected := range expectedArgs {
		if cmd.Args[i] != expected {
			t.Errorf("Expected arg[%d] '%s', got '%s'", i, expected, cmd.Args[i])
		}
	}

	if cmd.Stdout == nil {
		t.Error("Expected Stdout to be set")
	}

	if cmd.Stderr == nil {
		t.Error("Expected Stderr to be set")
	}

	if cmd.SysProcAttr == nil || !cmd.SysProcAttr.Setpgid {
		t.Error("Expected process group to be configured")
	}
}
