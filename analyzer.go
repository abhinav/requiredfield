package requiredfield

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

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
		Fset:             pass.Fset,
		Info:             pass.TypesInfo,
		ImportObjectFact: pass.ImportObjectFact,
		Reportf:          pass.Reportf,
	}
	e.Enforce(inspect)

	return nil, nil
}

type enforcer struct {
	Fset *token.FileSet
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
		var obj types.Object
		switch typ := typ.(type) {
		case *types.Named:
			// struct
			obj = typ.Obj()
		case *types.Pointer:
			// pointer to struct
			if named, ok := typ.Elem().(*types.Named); ok {
				obj = named.Obj()
			}
		}
		if obj == nil {
			return
		}

		var reqFields hasRequiredFields
		if !e.ImportObjectFact(obj, &reqFields) || len(reqFields.List) == 0 {
			return
		}

		unset := make(map[string]struct{})
		for _, f := range reqFields.List {
			unset[f] = struct{}{}
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
	new(ast.GenDecl),
	new(ast.TypeSpec),
	new(ast.StructType),
}

func (f *finder) Find(inspect *inspector.Inspector) {
	var curType *ast.Ident
	inspect.Nodes(_finderNodeFilter, func(n ast.Node, push bool) bool {
		switch n := n.(type) {
		case *ast.GenDecl:
			// Only look at type declarations.
			return n.Tok == token.TYPE

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

func (f *finder) structType(name *ast.Ident, t *ast.StructType) {
	var requiredFields []string
	for _, field := range t.Fields.List {
		if field.Comment == nil {
			continue
		}
		if !slices.ContainsFunc(field.Comment.List, isRequiredComment) {
			continue
		}

		for _, n := range field.Names {
			requiredFields = append(requiredFields, n.Name)
		}
	}

	if len(requiredFields) == 0 {
		return
	}
	slices.Sort(requiredFields)

	// Attach the fact to the type.
	obj, ok := f.Info.Defs[name]
	if !ok {
		f.Reportf(name.Pos(), "could not find object for %v", name)
		return
	}

	f.ExportObjectFact(obj, &hasRequiredFields{
		List: requiredFields,
	})
}

// hasRequiredFields is a Fact attached to structs.
// These fields must be explicitly set during initialization of the struct.
type hasRequiredFields struct {
	// List is a list of field names that are marked as required.
	List []string
}

var _ analysis.Fact = (*hasRequiredFields)(nil)

func (*hasRequiredFields) AFact() {}

func (f *hasRequiredFields) String() string {
	return "required<" + strings.Join(f.List, ", ") + ">"
}

func isRequiredComment(c *ast.Comment) bool {
	return c.Text == "// required"
}
