package a

import "fmt"

type AllOptional struct {
	X string
	Y int
}

type OneRequired struct { // want OneRequired:"required<B>"
	A string
	B int // required
	C bool
}

type twoRequired struct { // want twoRequired:"required<A, C>"
	A string // required
	B int
	C bool // required
	D float64
}

type allRequiredUnexported struct { // want allRequiredUnexported:"required<bar, baz, foo>"
	foo string // required
	bar int    // required
	baz bool   // required
}

type RequiredExported struct { // want RequiredExported:"required<A, B>"
	A string // required
	B int    // required
}

type aliasedStruct = RequiredExported

func x() {
	fmt.Println(AllOptional{})

	// Missing all fields:
	fmt.Println(OneRequired{}) // want "missing required fields: B"
	fmt.Println(twoRequired{}) // want "missing required fields: A, C"

	fmt.Println(allRequiredUnexported{}) // want "missing required fields: bar, baz, foo"

	// Missing some fields:
	fmt.Println(OneRequired{A: "a"}) // want "missing required fields: B"
	fmt.Println(twoRequired{A: "a"}) // want "missing required fields: C"

	fmt.Println(allRequiredUnexported{baz: false}) // want "missing required fields: bar, foo"

	fmt.Println("pointer:", &RequiredExported{}) // want "missing required fields: A, B"

	// Inside a map
	fmt.Println(map[string]RequiredExported{
		"foo": {},       // want "missing required fields: A, B"
		"bar": {A: "a"}, // want "missing required fields: B"
		"baz": {A: "a", B: 1},
	})

	// As a map key
	fmt.Println(map[RequiredExported]string{
		{A: "a", B: 1}: "foo",
		{A: "a"}:       "bar", // want "missing required fields: B"
		{B: 1}:         "baz", // want "missing required fields: A"
		{}:             "qux", // want "missing required fields: A, B"
	})

	// Aliased.
	fmt.Println(aliasedStruct{}) // want "missing required fields: A, B"

	// Unkeyed struct literals.
	fmt.Println(RequiredExported{}) // want "missing required fields: A, B"
	fmt.Println(RequiredExported{"a", 1})
}
