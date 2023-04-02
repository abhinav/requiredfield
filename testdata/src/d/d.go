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

func _1() {
	fmt.Println(Foo{}) // want "missing required fields: X, Z"
}

type Handler struct { // want Handler:"required<Callback>"
	Callback func(
		ctx context.Context,
		req *http.Request,
	) // required
}

func _2() {
	fmt.Println(Handler{}) // want "missing required fields: Callback"
}
