package main

import (
	"log/slog"
	"testing"
)

func TestIgnore(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		path        string
		ignorePaths []string
		want        bool
	}{
		{"top-level match", "vendor/x.go", []string{"vendor"}, true},
		{"directory itself", "vendor", []string{"vendor"}, true},
		{"nested match", "pkg/vendor/x.go", []string{"vendor"}, true},
		{"absolute match", "/repo/vendor/x.go", []string{"vendor"}, true},
		{"no match", "pkg/svc.go", []string{"vendor"}, false},
		{"prefix is not a full segment", "vendored/x.go", []string{"vendor"}, false},
		{"multiple paths: second matches", ".github/tools/t.go", []string{"vendor", ".github"}, true},
		{"multiple paths: none match", "cmd/main.go", []string{"vendor", ".github"}, false},
		{"empty ignore list", "vendor/x.go", nil, false},
		{"empty entries are skipped", "cmd/main.go", []string{"", ""}, false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ignore(tt.path, tt.ignorePaths); got != tt.want {
				t.Errorf("ignore(%q, %v) = %v, want %v", tt.path, tt.ignorePaths, got, tt.want)
			}
		})
	}
}

func TestExtractSATDs(t *testing.T) {
	NewLogger()

	// input
	workspacePath := "."
	ignorePaths := []string{"vendor"}

	// process & output
	satds, err := extractSATDs(workspacePath, ignorePaths)
	if err != nil {
		t.Fail()
	}

	// verify
	if len(satds) == 0 {
		t.Fail()
	}

	logger.Debug("satds", slog.Any("satds", satds))
}
