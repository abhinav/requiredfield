package requiredfield

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/ssa"
)

type enforcer struct {
	Fset *token.FileSet // required
	Info *types.Info    // required

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
		sort.Strings(missing)

		e.Reportf(lit.Lbrace, "missing required fields: %s", strings.Join(missing, ", "))
	})
}

func (e *enforcer) Enforce2(inspect *inspector.Inspector, funcs []*ssa.Function) {
	if len(funcs) == 0 {
		return
	}

	if filepath.Base(e.Fset.File(funcs[0].Pos()).Name()) != "e.go" {
		e.Enforce(inspect)
		return
	}

	fe := funcEnforcer{
		Fset:             e.Fset,
		ImportObjectFact: e.ImportObjectFact,
		Reportf:          e.Reportf,
	}
	for _, fn := range funcs {
		fe.Enforce(fn)
	}

	e.Enforce(inspect)
}

type structValueFact struct {
	Value ssa.Value
	Type  *types.Struct
	Block *ssa.BasicBlock
	Needs []string // list of required field names

	// Stores is a map from field name to a list of instructions
	// that set that field.
	// Multiple instructions may set the same field.
	Stores map[string][]ssa.Instruction

	// Uses is a list of instructions that use the struct.
	Uses []ssa.Instruction
}

type structStoreFact struct {
	Field string
	Inst  ssa.Instruction
	Block *ssa.BasicBlock
}

type structUseFact struct {
	Inst  ssa.Instruction
	Block *ssa.BasicBlock
}

type funcEnforcer struct {
	Fset             *token.FileSet                                       // required
	ImportObjectFact func(obj types.Object, fact analysis.Fact) bool      // required
	Reportf          func(pos token.Pos, msg string, args ...interface{}) // required

	structs map[ssa.Value]*structValueFact
}

func (e *funcEnforcer) Enforce(fn *ssa.Function) {
	// fn.WriteTo(os.Stderr)

	// TODO: Current method isn't right.
	// We need to go through the uses in-order.

	e.structs = make(map[ssa.Value]*structValueFact)
	for _, b := range fn.Blocks {
		for _, inst := range b.Instrs {
			e.visitInst(b, inst)
		}
	}

	// if fn.Name() == "check" {
	// 	fmt.Println("--- check")
	// 	fn.WriteTo(os.Stdout)
	// 	for _, p := range fn.Params {
	// 		fmt.Printf("param %#v\n", p)
	// 		for _, ref := range *p.Referrers() {
	// 			fmt.Printf("  ref %#v\n", ref)
	// 		}
	// 	}
	// }

	for _, s := range e.structs {
		for _, use := range s.Uses {
			e.checkStructUse(s, use)
		}
	}
}

func (e *funcEnforcer) checkStructUse(s *structValueFact, use ssa.Instruction) {
	useBlock := use.Block()

	var missing []string
fields:
	for _, name := range s.Needs {
		if sets, ok := s.Stores[name]; ok {
			// Check that at least one of the stores dominates the use.
			for _, set := range sets {
				if set.Block().Dominates(useBlock) {
					continue fields
				}
			}
		}
		missing = append(missing, name)
	}

	if len(missing) == 0 {
		return
	}

	sort.Strings(missing)
	e.Reportf(use.Pos(), "missing required fields: %s", strings.Join(missing, ", "))
	// fmt.Printf("%v:missing required fields: %s\n", e.Fset.Position(use.Pos()), strings.Join(missing, ", "))
}

func (e *funcEnforcer) visitInst(b *ssa.BasicBlock, inst ssa.Instruction) {
	switch inst := inst.(type) {
	case *ssa.Alloc:
		e.visitAlloc(b, inst)

	case *ssa.FieldAddr:
		// Ignore

	case *ssa.Store:
		e.visitStore(inst)

	// All other instructions are checked for uses of
	// known structs.
	default:
		// TODO:
		// are there instructions that use the struct
		// that we should ignore?
		//
		// TODO: pointers to structs
		for _, o := range inst.Operands(nil) {
			e.visitValue(inst, *o)
		}
	}
}

func (e *funcEnforcer) visitAlloc(b *ssa.BasicBlock, inst *ssa.Alloc) {
	var needsFields []string
	typ := inst.Type()
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	var st *types.Struct
	switch typ := typ.(type) {
	case *types.Named:
		var reqFields hasRequiredFields
		if !e.ImportObjectFact(typ.Obj(), &reqFields) {
			// not a struct or does not have
			// required fields.
			return
		}
		st = typ.Underlying().(*types.Struct) // TODO: check
		needsFields = reqFields.List

	case *types.Struct:
		// Anonymous struct.
		var fact isRequiredField
		for i := 0; i < typ.NumFields(); i++ {
			f := typ.Field(i)
			if e.ImportObjectFact(f, &fact) {
				needsFields = append(needsFields, f.Name())
			}
		}
		st = typ
		// TODO: field indexes might be easier?

	default:
		// Not a struct or does not have required fields.
		return
	}

	// Not a struct or does not have required fields.
	if len(needsFields) == 0 {
		return
	}

	sort.Strings(needsFields)

	// fmt.Printf("%v:declare struct %v\n", e.Fset.Position(inst.Pos()), inst)
	e.structs[inst] = &structValueFact{
		Value:  inst,
		Type:   st,
		Block:  b,
		Needs:  needsFields,
		Stores: make(map[string][]ssa.Instruction),
	}
}

func (e *funcEnforcer) visitStore(inst *ssa.Store) {
	// Check if we're storing a value into a struct field.
	switch dst := inst.Addr.(type) {
	case *ssa.Alloc:
		// If we're just recording a function parameter,
		// forget about the struct.
		if _, ok := inst.Val.(*ssa.Parameter); ok {
			delete(e.structs, dst)
		}

	case *ssa.FieldAddr:
		fact, ok := e.structs[dst.X]
		if !ok {
			return
		}

		// TODO We can do better.
		fname := fact.Type.Field(dst.Field).Name()
		fact.Stores[fname] = append(fact.Stores[fname], inst)
	}

	e.visitValue(inst, inst.Val)
}

func (e *funcEnforcer) visitValue(inst ssa.Instruction, v ssa.Value) {
	fact, ok := e.structs[v]
	if !ok {
		return
	}

	// TODO: Don't bother storing the use
	// if there's already a use that dominates it.
	fact.Uses = append(fact.Uses, inst)
}
