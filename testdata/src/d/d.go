package d

import (
	"context"
	"fmt"
	"net/http"
)

type Foo struct { // want Foo:"required<X, Z>"
	X struct {
		// multi-line type
	} // required

	// required: should be ignored
	Y string
	// required: should be ignored

	Z int `json:"z"` // required: has field tag
}

func _() {
	fmt.Println(Foo{}) // want "missing required fields: X, Z"
}

type Handler struct { // want Handler:"required<Callback>"
	Callback func(
		ctx context.Context,
		req *http.Request,
	) // required
}

func _() {
	fmt.Println(Handler{}) // want "missing required fields: Callback"
}

type irregularSpacing struct { // want irregularSpacing:"required<A, B, C, D, E>"
	A int // required
	B int //  required
	C int //     required
	D int // required: some context
	E int //   required: some context
	F int /* required: not enforced */
}

func _() {
	fmt.Println(irregularSpacing{}) // want "missing required fields: A, B, C, D, E"
}

type invalidTags struct { // want invalidTags:"required<A>"
	A int // required
	B int // requiredsuffix
	C int // prefixrequired
}

func _() {
	fmt.Println(invalidTags{}) // want "missing required fields: A"
}
