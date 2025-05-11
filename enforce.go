package requiredfield

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

type enforcer struct {
	Info *types.Info // required

	ImportObjectFact func(obj types.Object, fact analysis.Fact) bool      // required
	Reportf          func(pos token.Pos, msg string, args ...interface{}) // required
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
		if alias, ok := typ.(*types.Alias); ok {
			typ = types.Unalias(alias)
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
				// Elements will not be KeyValueExprs
				// only if unkeyed struct literal is used.
				// In case of unkeyed literals,
				// the compiler enforces that
				// all fields are specified,
				// so there's nothing for us to do
				// for this struct.
				return
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
		sort.Strings(missing)

		e.Reportf(lit.Lbrace, "missing required fields: %s", strings.Join(missing, ", "))
	})
}
