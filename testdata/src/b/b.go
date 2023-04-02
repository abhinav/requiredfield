package b

import (
	"a"
	"fmt"
)

type MyStruct struct { // want MyStruct:"required<MyOwn>"
	Thing a.RequiredExported

	MyOwn string // required
}

type User struct { // want User:"required<ID>"
	// ID is the user's unique identifier.
	ID string // required

	// DisplayName is an optional display name.
	DisplayName string // optional
}

func _1() {
	fmt.Println(a.AllOptional{})

	// Missing all fields:
	fmt.Println(a.OneRequired{})      // want "missing required fields: B"
	fmt.Println(a.RequiredExported{}) // want "missing required fields: A, B"

	// Missing some fields:
	fmt.Println(a.OneRequired{A: "a"})     // want "missing required fields: B"
	fmt.Println(a.RequiredExported{A: ""}) // want "missing required fields: B"

	// Missing no required fields:
	fmt.Println(a.OneRequired{B: 1})
	fmt.Println(a.RequiredExported{A: "", B: 1})

	fmt.Println("pointer:", &a.RequiredExported{}) // want "missing required fields: A, B"

	// Inside another struct:
	fmt.Println(MyStruct{ // want "missing required fields: MyOwn"
		Thing: a.RequiredExported{}, // want "missing required fields: A, B"
	})
	fmt.Println(MyStruct{
		MyOwn: "foo",
		// Can omit Thing because it's not required.
	})
	fmt.Println(MyStruct{
		MyOwn: "foo",
		Thing: a.RequiredExported{ // want "missing required fields: A"
			B: 1,
		},
	})

	// Inside a slice
	fmt.Println([]User{
		{ // want "missing required fields: ID"
			DisplayName: "foo",
		},
		{ID: "foo"},
		{ID: "foo", DisplayName: "bar"},
	})
}
