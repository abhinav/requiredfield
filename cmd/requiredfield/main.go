// requiredfield is a linter for Go code
// that checks for missing required fields in struct literals.
// See top-level package documentation for more information on the linter.
//
// # Usage
//
// To use this linter, install the requiredfield binary:
//
//	$ go install go.abhg.dev/requiredfield/cmd/requiredfield@latest
//
// Pass the path to this file to 'go vet':
//
//	$ go vet -vettool=$(which requiredfield) ./...
//
// # As a golangci-lint plugin
//
// Build the plugin:
//
//	$ go build -buildmode=plugin go.abhg.dev/requiredfield/cmd/requiredfield
//
// Then enable it in the golangci-lint configuration:
//
//	$ cat .golangci.yml
//	linter-settings:
//	  custom:
//	    requiredfield:
//	      path: requiredfield.so
//	      description: Checks for required fields in structs
//	      original-url: go.abhg.dev/requiredfield
package main

import (
	"go.abhg.dev/requiredfield"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(requiredfield.Analyzer)
}

// AnalyzerPlugin provides the analyzer as a golangci-lint plugin.
var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

// GetAnalyzers returns the requiredfield analyzer.
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{requiredfield.Analyzer}
}
