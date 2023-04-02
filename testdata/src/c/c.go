package c

import (
	"a"
	"fmt"
)

type Foo struct { // want Foo:"required<OneRequired>"
	a.OneRequired // required
	a.RequiredExported

	Bar struct {
		X int // required // want X:"required"
	}
}

func _1() {
	fmt.Println(Foo{}) // want "missing required fields: OneRequired"

	fmt.Println(Foo{
		OneRequired: a.OneRequired{}, // want "missing required fields: B"
		RequiredExported: a.RequiredExported{ // want "missing required fields: A"
			B: 42,
		},
	})
}

type Bar map[string]struct {
	X int // required // want X:"required"
}

func _2() {
	fmt.Println(Bar{
		"foo": {}, // want "missing required fields: X"
	})
}

type Alias = a.OneRequired

func _3() {
	fmt.Println(Alias{}) // want "missing required fields: B"
}

type embedsAlias struct { // want embedsAlias:"required<Alias>"
	Alias // required
}

func _4() {
	fmt.Println(embedsAlias{}) // want "missing required fields: Alias"

	fmt.Println(embedsAlias{
		Alias: Alias{}, // want "missing required fields: B"
	})
}

type embedsPtr struct { // want embedsPtr:"required<OneRequired>"
	*a.OneRequired // required
}

type embedsAliasPtr struct { // want embedsAliasPtr:"required<Alias>"
	*Alias // required
}

func _5() {
	fmt.Println(embedsPtr{})      // want "missing required fields: OneRequired"
	fmt.Println(embedsAliasPtr{}) // want "missing required fields: Alias"
}
