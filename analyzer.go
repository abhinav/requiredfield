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
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

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

type enforcer struct {
	Info *types.Info

	ImportObjectFact func(obj types.Object, fact analysis.Fact) bool
	Reportf          func(pos token.Pos, msg string, args ...interface{})
}

var _enforceNodeFilter = []ast.Node{
	new(ast.CompositeLit),
}

func (e *enforcer) Enforce(inspect *inspector.Inspector) {
	inspect.Preorder(_enforceNodeFilter, func(n ast.Node) {
		lit := n.(*ast.CompositeLit)
		typ := e.Info.TypeOf(lit)

		if ptr, ok := typ.(*types.Pointer); ok {
			typ = ptr.Elem()
		}

		var unset map[string]struct{}
		switch typ := typ.(type) {
		case *types.Named:
			// named struct (probably)
			var reqFields hasRequiredFields
			if !e.ImportObjectFact(typ.Obj(), &reqFields) {
				return
			}

			unset = make(map[string]struct{}, len(reqFields.List))
			for _, name := range reqFields.List {
				unset[name] = struct{}{}
			}
		case *types.Struct:
			// anonymous struct
			var fact isRequiredField
			for i := 0; i < typ.NumFields(); i++ {
				f := typ.Field(i)
				if e.ImportObjectFact(f, &fact) {
					if unset == nil {
						unset = make(map[string]struct{})
					}
					unset[f.Name()] = struct{}{}
				}
			}
		}

		if len(unset) == 0 {
			return
		}

		// Check that all required fields are set.
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			id, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			delete(unset, id.Name)
		}

		if len(unset) == 0 {
			return
		}

		var missing []string
		for f := range unset {
			missing = append(missing, f)
		}
		slices.Sort(missing)

		e.Reportf(lit.Lbrace, "missing required fields: %s", strings.Join(missing, ", "))
	})
}

type finder struct {
	Info *types.Info

	ExportObjectFact func(obj types.Object, fact analysis.Fact)
	Reportf          func(pos token.Pos, msg string, args ...interface{})
}

var _finderNodeFilter = []ast.Node{
	new(ast.TypeSpec),
	new(ast.StructType),
}

func (f *finder) Find(inspect *inspector.Inspector) {
	var curType *ast.Ident
	inspect.Nodes(_finderNodeFilter, func(n ast.Node, push bool) bool {
		switch n := n.(type) {
		case *ast.TypeSpec:
			if push {
				curType = n.Name
			} else {
				curType = nil
			}

		case *ast.StructType:
			if push {
				f.structType(curType, n)
			}
		}

		return true
	})
}

// structType inspects the provided struct definition.
// If it has any required fields, it attaches a fact to the type.
// name may be nil if the struct is anonymous.
func (f *finder) structType(name *ast.Ident, t *ast.StructType) {
	var (
		requiredIndexes []int
		requiredFields  []string
	)
	for i, field := range t.Fields.List {
		if field.Comment == nil {
			continue
		}
		if !slices.ContainsFunc(field.Comment.List, isRequiredComment) {
			continue
		}

		requiredIndexes = append(requiredIndexes, i)
		for _, n := range field.Names {
			requiredFields = append(requiredFields, n.Name)
		}
	}

	if len(requiredFields) == 0 {
		return
	}
	slices.Sort(requiredFields)

	if name != nil {
		// Named struct.
		// Attach the fact to the type.
		obj, ok := f.Info.Defs[name]
		if !ok {
			f.Reportf(name.Pos(), "could not find object for %v", name)
			return
		}
		f.ExportObjectFact(obj, &hasRequiredFields{
			List: requiredFields,
		})
	} else {
		// Anonymous struct.
		// Attach to individual fields.
		st, ok := f.Info.TypeOf(t).(*types.Struct)
		if !ok {
			return
		}

		for _, i := range requiredIndexes {
			f.ExportObjectFact(st.Field(i), &isRequiredField{})
		}
	}
}

// hasRequiredFields is a Fact attached to structs
// listing its required fields.
type hasRequiredFields struct {
	// List is a list of field names
	// in the struct that are marked required.
	List []string
}

var _ analysis.Fact = (*hasRequiredFields)(nil)

func (*hasRequiredFields) AFact() {}

func (f *hasRequiredFields) String() string {
	return "required<" + strings.Join(f.List, ", ") + ">"
}

// isRequiredField is a Fact attached to fields of anonymous structs
// that are marked required.
type isRequiredField struct{}

var _ analysis.Fact = (*isRequiredField)(nil)

func (*isRequiredField) AFact() {}

func (f *isRequiredField) String() string {
	return "required"
}

const _required = "// required"

func isRequiredComment(c *ast.Comment) bool {
	if c.Text == _required {
		return true
	}

	if !strings.HasPrefix(c.Text, _required) {
		return false
	}

	for _, r := range c.Text[len(_required):] {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
	}

	// Impossible:
	// If the comment is not "// required", but it starts with that,
	// the loop above will always return before we get here.
	// This is just to make the compiler happy.
	return false
}
