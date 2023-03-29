package requiredfield

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

type finder struct {
	Fset *token.FileSet // required
	Info *types.Info    // required

	ExportObjectFact func(obj types.Object, fact analysis.Fact)           // required
	Reportf          func(pos token.Pos, msg string, args ...interface{}) // required
}

var _finderNodeFilter = []ast.Node{
	new(ast.TypeSpec),
	new(ast.StructType),
}

func (f *finder) Find(inspect *inspector.Inspector) {
	seen := make(map[*ast.StructType]struct{})

	inspect.Preorder(_finderNodeFilter, func(n ast.Node) {
		var (
			name *ast.Ident
			st   *ast.StructType
		)

		switch n := n.(type) {
		case *ast.TypeSpec:
			if t, ok := n.Type.(*ast.StructType); ok {
				name = n.Name
				st = t
			}
		case *ast.StructType:
			st = n
		}

		// If the type spec is not a struct, or if we've already seen it,
		// we can skip it.
		if st == nil {
			return
		}
		if _, ok := seen[st]; ok {
			return
		}

		seen[st] = struct{}{}
		f.structType(name, st)
	})
}

// structType inspects the provided struct definition.
// If it has any required fields, it attaches a fact to the type.
// name may be nil if the struct is anonymous.
func (f *finder) structType(name *ast.Ident, t *ast.StructType) {
	file := f.Fset.File(t.Pos())

	var (
		requiredIndexes []int
		requiredFields  []string
	)
	st := f.Info.TypeOf(t).(*types.Struct)
	for i, field := range t.Fields.List {
		if field.Comment == nil {
			continue
		}

		var found bool
		fieldLine := file.Line(field.End())
		for _, c := range field.Comment.List {
			if file.Line(c.Pos()) == fieldLine && isRequiredComment(c) {
				found = true
				break
			}
		}

		if !found {
			continue
		}

		requiredIndexes = append(requiredIndexes, i)
		if field.Names == nil {
			// Embedded fields don't have field.Names.
			name := st.Field(i).Name()
			requiredFields = append(requiredFields, name)
		} else {
			for _, n := range field.Names {
				requiredFields = append(requiredFields, n.Name)
			}
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
