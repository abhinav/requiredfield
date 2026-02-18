package h

import (
	"fmt"
	"g"
)

// Non-generic alias to an instantiated generic struct.
type IntContainer = g.Container[int]

// Parameterized alias forwarding type params.
type Alias[T any] = g.Container[T]

// Parameterized alias with narrower constraints.
type NumContainer[T interface{ ~int | ~float64 }] = g.Container[T]

// Chained aliases: parameterized alias to another alias.
type AliasOfAlias[T any] = Alias[T]

// Chained alias: non-generic alias to a parameterized alias.
type StringAlias = Alias[string]

// Double chain: non-generic alias to chained alias.
type DoubleChain = AliasOfAlias[bool]

func nonGenericAliasToGeneric() {
	fmt.Println(IntContainer{})           // want "missing required fields: Value"
	fmt.Println(IntContainer{Label: "x"}) // want "missing required fields: Value"
	fmt.Println(IntContainer{Value: 42})
	fmt.Println(IntContainer{Value: 0, Label: "zero"})
}

func parameterizedAlias() {
	fmt.Println(Alias[int]{})            // want "missing required fields: Value"
	fmt.Println(Alias[string]{})         // want "missing required fields: Value"
	fmt.Println(Alias[int]{Value: 1})
	fmt.Println(Alias[string]{Value: "x", Label: "ok"})
}

func narrowerConstraintAlias() {
	fmt.Println(NumContainer[int]{})     // want "missing required fields: Value"
	fmt.Println(NumContainer[float64]{}) // want "missing required fields: Value"
	fmt.Println(NumContainer[int]{Value: 42})
}

func chainedAliases() {
	fmt.Println(AliasOfAlias[int]{})    // want "missing required fields: Value"
	fmt.Println(AliasOfAlias[int]{Value: 1})

	fmt.Println(StringAlias{})          // want "missing required fields: Value"
	fmt.Println(StringAlias{Value: "x"})

	fmt.Println(DoubleChain{})          // want "missing required fields: Value"
	fmt.Println(DoubleChain{Value: true})
}

func pointerToAliasedGeneric() {
	fmt.Println(&IntContainer{})    // want "missing required fields: Value"
	fmt.Println(&Alias[int]{})      // want "missing required fields: Value"
	fmt.Println(&IntContainer{Value: 1})
	fmt.Println(&Alias[int]{Value: 1})
}

// Struct that embeds a generic type (marked required).
type EmbedGeneric struct { // want EmbedGeneric:"required<Container>"
	g.Container[int] // required
}

// Struct that embeds a generic alias (marked required).
type EmbedAlias struct { // want EmbedAlias:"required<IntContainer>"
	IntContainer // required
}

// Struct that embeds a pointer to a generic type.
type EmbedGenericPtr struct { // want EmbedGenericPtr:"required<Container>"
	*g.Container[int] // required
}

func embeddingGenericTypes() {
	fmt.Println(EmbedGeneric{})    // want "missing required fields: Container"
	fmt.Println(EmbedGeneric{
		Container: g.Container[int]{}, // want "missing required fields: Value"
	})
	fmt.Println(EmbedGeneric{
		Container: g.Container[int]{Value: 1},
	})

	fmt.Println(EmbedAlias{})      // want "missing required fields: IntContainer"
	fmt.Println(EmbedAlias{
		IntContainer: IntContainer{}, // want "missing required fields: Value"
	})
	fmt.Println(EmbedAlias{
		IntContainer: IntContainer{Value: 1},
	})

	fmt.Println(EmbedGenericPtr{}) // want "missing required fields: Container"
}

func genericStructInMapViaAlias() {
	fmt.Println(map[string]IntContainer{
		"a": {},          // want "missing required fields: Value"
		"b": {Value: 42},
	})

	fmt.Println(map[string]Alias[int]{
		"x": {},         // want "missing required fields: Value"
		"y": {Value: 1},
	})
}

// Pair alias to verify multi-field generics through aliases.
type IntPair = g.Pair[string, int]

func pairAlias() {
	fmt.Println(IntPair{})         // want "missing required fields: Key, Value"
	fmt.Println(IntPair{Key: "x"}) // want "missing required fields: Value"
	fmt.Println(IntPair{Key: "x", Value: 42})
}

// Cross-package generic usage without aliases.
func crossPackageDirect() {
	fmt.Println(g.Container[int]{})    // want "missing required fields: Value"
	fmt.Println(g.Container[int]{Value: 1})
	fmt.Println(g.Pair[string, int]{}) // want "missing required fields: Key, Value"
	fmt.Println(g.Pair[string, int]{Key: "a", Value: 1})
}
