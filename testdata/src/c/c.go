package c

import (
	"a"
	"fmt"
)

type Foo struct {
	a.OneRequired
	a.RequiredExported

	Bar struct {
		X int // required // want X:"required"
	}
}

func _() {
	f := Foo{
		OneRequired: a.OneRequired{}, // want "missing required fields: B"
		RequiredExported: a.RequiredExported{ // want "missing required fields: B"
			A: "",
		},
	}
	fmt.Println(f)
}

type Bar map[string]struct {
	X int // required // want X:"required"
}

func _() {
	fmt.Println(Bar{
		"foo": {}, // want "missing required fields: X"
	})
}
