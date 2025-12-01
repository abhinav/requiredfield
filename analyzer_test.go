package requiredfield

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testDataDir := analysistest.TestData()
	srcDir := filepath.Join(testDataDir, "src")

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("failed to read src directory: %v", err)
	}

	// Collect all test packages.
	var packages []string
	for _, entry := range entries {
		if entry.IsDir() {
			packages = append(packages, entry.Name())
		}
	}

	for _, pkg := range packages {
		t.Run(pkg, func(t *testing.T) {
			t.Parallel()

			var linter requiredfieldLinter

			rcPath := filepath.Join(srcDir, pkg, "requiredfield.rc")
			if bs, err := os.ReadFile(rcPath); err == nil {
				err := linter.Config.Parse(strings.NewReader(string(bs)))
				if err != nil {
					t.Fatalf("failed to parse requiredfield.rc: %v", err)
				}
			}

			analysistest.Run(t, testDataDir, linter.Analyzer(), pkg)
		})
	}
}
