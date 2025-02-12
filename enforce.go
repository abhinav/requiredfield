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
	inspect.WithStack(_enforceNodeFilter, func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		if !push {
			return true
		}

		e.visit(n, stack)
		return true
	})
}

func (e *enforcer) visit(n ast.Node, stack []ast.Node) {
	lit := n.(*ast.CompositeLit)
	typ := e.Info.TypeOf(lit)

	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}
	if alias, ok := typ.(*types.Alias); ok {
		typ = types.Unalias(alias)
	}

	var unset map[string]struct{} // required fields that are not set
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
		// Type has no required fields, or is not a struct.
		return
	}

	// If a function returns a struct and an error,
	// it's common to do 'return MyStruct{}, err' for failures.
	// It's not desirable to enforce required fields in this case.
	//
	// Therefore, if we encounter a struct literal in a return,
	// and the last value of that return statement is an error
	// that is not explicitly set to "nil",
	// we will not enforce required fields on that struct literal.
	if len(stack) > 1 && e.isReturnedWithNonNilError(stack) {
		// The struct literal is part of a return statement
		// that has a non-nil error as its last return value.
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
}

// isReturnedWithNonNilError reports whether target is part of a return
// statement that has a non-nil error as its last return value,
// but is not itself the last return value or a subexpression of it.
//
// The "but" is important:
// if the target is the last return value (or part of it),
// it should still be checked for required fields, e.g.
//
//	return nil, &MyError{...} // and MyError has required fields
func (e *enforcer) isReturnedWithNonNilError(stack []ast.Node) bool {
	// Find the nearest return statement.
	var retStmt *ast.ReturnStmt
	retIdx := -1
	for idx := len(stack) - 1; idx >= 0; idx-- {
		var ok bool
		if retStmt, ok = stack[idx].(*ast.ReturnStmt); ok {
			retIdx = idx
			break
		}
	}
	if retIdx == -1 {
		// No return statement found.
		return false
	}

	// Find the nearest function's type for the return statement.
	var ftype *ast.FuncType
	for idx := retIdx - 1; idx >= 0; idx-- {
		switch n := stack[idx].(type) {
		case *ast.FuncDecl:
			ftype = n.Type
		case *ast.FuncLit:
			ftype = n.Type
		}
	}
	if ftype == nil {
		// Impossible, but we don't want to panic.
		return false
	}

	// The last return type must be an error.
	returnTypes := ftype.Results.List
	if len(returnTypes) == 0 {
		// No return types.
		return false
	}
	lastReturnType, ok := returnTypes[len(returnTypes)-1].Type.(*ast.Ident)
	if !ok || lastReturnType.Name != "error" {
		// The last return type is not "error".
		return false
	}

	// If last return value is nil, we want to enforce required fields.
	lastReturn := retStmt.Results[len(retStmt.Results)-1]
	if id, ok := lastReturn.(*ast.Ident); ok && id.Name == "nil" {
		return false
	}

	// At this point, we know this is a return statement in a function
	// where the last return type is an error,
	// and the value is definitely not the nil identifier.
	//
	// We want to ignore this node (return true) only if
	// the target is not part of the last return value itself.
	//
	// To do this, we'll verify that lastReturn is not equal to,
	// or an ancestor of the target. (stack[len(stack)-1] is target)
	for idx := len(stack) - 1; idx > retIdx; idx-- {
		if stack[idx] == lastReturn {
			// Target is part of the last return value.
			return false
		}
	}

	// Target is not part of the last return value.
	return true
}
