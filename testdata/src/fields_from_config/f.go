package f

import (
	"external"
	"fmt"
)

// External fields are marked from the requiredfield.rc file.

func useExternal() {
	fmt.Println(external.User{})              // want "missing required fields: ID, Name"
	fmt.Println(external.User{ID: "123"})     // want "missing required fields: Name"
	fmt.Println(external.User{Name: "Alice"}) // want "missing required fields: ID"
	fmt.Println(external.User{ID: "123", Name: "Alice"})

	fmt.Println(external.Config{})            // want "missing required fields: APIKey"
	fmt.Println(external.Config{Timeout: 10}) // want "missing required fields: APIKey"
	fmt.Println(external.Config{APIKey: "secret"})
}

// Some fields of LocalType are also marked required.
// Both get merged.

type LocalType struct { // want LocalType:"required<A>"
	A string // required
	B int
	C bool
}

func useLocal() {
	fmt.Println(LocalType{})           // want "missing required fields: A, B"
	fmt.Println(LocalType{A: "value"}) // want "missing required fields: B"
	fmt.Println(LocalType{B: 42})      // want "missing required fields: A"
	fmt.Println(LocalType{A: "value", B: 42})
}
