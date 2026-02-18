//go:build go1.26

package g

import "fmt"

// Self-referencing generic type (valid since Go 1.26).
type Node[N Node[N]] struct { // want Node:"required<ID>"
	ID       int // required
	Children []N
}

func selfReferencing() {
	fmt.Println(ConcreteNode{})    // want "missing required fields: ID"
	fmt.Println(ConcreteNode{ID: 1})
}
