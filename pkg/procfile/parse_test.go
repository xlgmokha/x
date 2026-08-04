package procfile

import (
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("TestParseBasic", func(t *testing.T) {
		input := "web: echo 'Hello World'\n"
		procs := Parse(strings.NewReader(input))

		if len(procs) != 1 {
			t.Errorf("Expected 1 process, got %d", len(procs))
		}

		proc := procs[0]
		if proc.name != "web" {
			t.Errorf("Expected name 'web', got '%s'", proc.name)
		}

		expectedArgs := []string{"echo", "'Hello", "World'"}
		args := proc.args
		if len(args) != len(expectedArgs) {
			t.Errorf("Expected %d args, got %d", len(expectedArgs), len(args))
		}

		for i, expected := range expectedArgs {
			if args[i] != expected {
				t.Errorf("Expected arg[%d] '%s', got '%s'", i, expected, args[i])
			}
		}
	})

	t.Run("TestParseMultiple", func(t *testing.T) {
		input := `web: ./server -port 8080
worker: ./worker -queue tasks
redis: redis-server --port 6379`

		procs := Parse(strings.NewReader(input))

		if len(procs) != 3 {
			t.Errorf("Expected 3 processes, got %d", len(procs))
		}

		expectedNames := []string{"web", "worker", "redis"}
		for i, expectedName := range expectedNames {
			if procs[i].name != expectedName {
				t.Errorf("Expected process[%d] name '%s', got '%s'", i, expectedName, procs[i].name)
			}
		}
	})

	t.Run("TestParseCommentsAndEmptyLines", func(t *testing.T) {
		input := `# This is a comment
web: echo 'Web server'

# Another comment
worker: echo 'Background worker'

# Empty lines should be ignored`

		procs := Parse(strings.NewReader(input))

		if len(procs) != 2 {
			t.Errorf("Expected 2 processes, got %d", len(procs))
		}

		if procs[0].name != "web" || procs[1].name != "worker" {
			t.Errorf("Expected processes 'web' and 'worker', got '%s' and '%s'",
				procs[0].name, procs[1].name)
		}
	})

	t.Run("TestParseInvalidLines", func(t *testing.T) {
		input := `web echo 'missing colon'
worker: echo 'this one is valid'
invalid line without colon`

		procs := Parse(strings.NewReader(input))

		if len(procs) != 1 {
			t.Errorf("Expected 1 process, got %d", len(procs))
		}

		if procs[0].name != "worker" {
			t.Errorf("Expected process name 'worker', got '%s'", procs[0].name)
		}
	})

	t.Run("TestParseEmptyCommands", func(t *testing.T) {
		input := `web: 
worker: echo 'this is valid'
empty:`

		procs := Parse(strings.NewReader(input))

		if len(procs) != 1 {
			t.Errorf("Expected 1 process, got %d", len(procs))
		}

		if procs[0].name != "worker" {
			t.Errorf("Expected process name 'worker', got '%s'", procs[0].name)
		}
	})

	t.Run("TestParseEnvironmentVariables", func(t *testing.T) {
		os.Setenv("TEST_PORT", "3000")
		os.Setenv("TEST_ENV", "development")
		defer func() {
			os.Unsetenv("TEST_PORT")
			os.Unsetenv("TEST_ENV")
		}()

		input := `web: ./server -port $TEST_PORT
worker: ./worker -env $TEST_ENV`

		procs := Parse(strings.NewReader(input))

		if len(procs) != 2 {
			t.Errorf("Expected 2 processes, got %d", len(procs))
		}

		// Check that environment variables were expanded
		webArgs := procs[0].args
		if len(webArgs) < 3 || webArgs[2] != "3000" {
			t.Errorf("Expected web port '3000', got args: %v", webArgs)
		}

		workerArgs := procs[1].args
		if len(workerArgs) < 3 || workerArgs[2] != "development" {
			t.Errorf("Expected worker env 'development', got args: %v", workerArgs)
		}
	})
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		file     string
		expected int
	}{
		{"testdata/valid/basic.procfile", 1},
		{"testdata/valid/multiple.procfile", 3},
		{"testdata/valid/with-comments.procfile", 2},
		{"testdata/invalid/no-colon.procfile", 1},      // Only valid lines parsed
		{"testdata/invalid/empty-command.procfile", 1}, // Only valid lines parsed
		{"testdata/env/variables.procfile", 3},         // Environment variable expansion
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			procs, err := ParseFile(test.file)
			if err != nil {
				t.Fatalf("ParseFile failed: %v", err)
			}

			if len(procs) != test.expected {
				t.Errorf("Expected %d processes, got %d", test.expected, len(procs))
			}
		})
	}
}
