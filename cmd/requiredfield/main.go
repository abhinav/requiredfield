// requiredfield is a linter for Go code
// that checks for missing required fields in struct literals.
// See top-level package documentation for more information on the linter.
//
// To use this linter, run the 'requiredfield' binary directly:
//
//	$ requiredfield ./...
//
// Alternatively, you can use the 'go vet' command:
//
//	$ go vet -vettool=$(which requiredfield) ./...
package main

import (
	"go.abhg.dev/requiredfield"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(requiredfield.Analyzer)
}
