package c

import (
	"a"
	"fmt"
)

type Foo struct {
	a.OneRequired
	a.RequiredExported
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
