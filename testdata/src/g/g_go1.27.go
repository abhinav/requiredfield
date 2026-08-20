//go:build go1.27

package g

import "fmt"

type Factory struct{}

func (Factory) New[T any]() Container[T] {
	return Container[T]{} // want "missing required fields: Value"
}

type Metadata struct {
	ID int
}

type Record struct { // want Record:"required<Metadata>"
	Metadata // required
}

func promotedFieldKey() {
	fmt.Println(Record{ID: 42})
	fmt.Println(Record{}) // want "missing required fields: Metadata"
}
