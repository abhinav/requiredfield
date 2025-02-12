//go:build go1.24

package e

import "fmt"

type AllOptional struct {
	X string
	Y int
}

type OneRequired struct { // want OneRequired:"required<B>"
	A string
	B int // required
}

type genericAlias[T any] = T

func _() {
	fmt.Println(genericAlias[AllOptional]{})
	fmt.Println(genericAlias[OneRequired]{}) // want "missing required fields: B"
	fmt.Println(genericAlias[OneRequired]{B: 42})
}
