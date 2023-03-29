// Package requiredfield implements a linter
// that checks for required fields during struct initialization.
//
// Fields can be marked as required by adding a comment in the form
// "// required" next to the field, optionally followed by a description.
// For example:
//
//	type T struct {
//		A string // required
//		B int    // required: must be positive
//		C bool   // required because reasons
//	}
//
// The analyzer will report an error when an instance of the struct is
// initialized without setting one or more of the required fields explicitly.
// For example:
//
//	T{A: "foo"} // error: missing required fields: B, C
//
// The explicit value can be the zero value of the field type,
// but it must be set explicitly.
//
//	T{A: "foo", B: 0, C: false}
package requiredfield

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer implements the requiredfield linter.
//
// See package documentation for details.
var Analyzer = &analysis.Analyzer{
	Name: "requiredfield",
	Doc:  "check for required fields during struct initialization",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	FactTypes: []analysis.Fact{
		new(isRequiredField),
		new(hasRequiredFields),
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	f := finder{
		Fset:             pass.Fset,
		Info:             pass.TypesInfo,
		ExportObjectFact: pass.ExportObjectFact,
		Reportf:          pass.Reportf,
	}
	f.Find(inspect)

	e := enforcer{
		Info:             pass.TypesInfo,
		ImportObjectFact: pass.ImportObjectFact,
		Reportf:          pass.Reportf,
	}
	e.Enforce(inspect)

	return nil, nil
}
