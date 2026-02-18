package g

import (
	"errors"
	"fmt"
	"math/rand"
)

// Generic struct with a required field.
type Container[T any] struct { // want Container:"required<Value>"
	Value T      // required
	Label string
}

// Generic struct with multiple required fields.
type Pair[K comparable, V any] struct { // want Pair:"required<Key, Value>"
	Key   K // required
	Value V // required
}

// Generic struct with all optional fields (no diagnostics expected).
type Box[T any] struct {
	Contents T
	Size     int
}

func basicGeneric() {
	// Missing required field.
	fmt.Println(Container[int]{})           // want "missing required fields: Value"
	fmt.Println(Container[string]{})        // want "missing required fields: Value"
	fmt.Println(Container[[]byte]{})        // want "missing required fields: Value"
	fmt.Println(Container[int]{Label: "x"}) // want "missing required fields: Value"

	// All required fields set.
	fmt.Println(Container[int]{Value: 42})
	fmt.Println(Container[int]{Value: 0, Label: "zero"})
	fmt.Println(Container[string]{Value: ""})

	// All optional struct is fine empty.
	fmt.Println(Box[int]{})
}

func multipleRequiredFields() {
	fmt.Println(Pair[string, int]{})             // want "missing required fields: Key, Value"
	fmt.Println(Pair[string, int]{Key: "x"})     // want "missing required fields: Value"
	fmt.Println(Pair[string, int]{Value: 42})     // want "missing required fields: Key"
	fmt.Println(Pair[string, int]{Key: "x", Value: 42})
}

func pointerToGenericStruct() {
	fmt.Println(&Container[int]{})    // want "missing required fields: Value"
	fmt.Println(&Container[int]{Value: 1})
	fmt.Println(&Pair[int, int]{})    // want "missing required fields: Key, Value"
	fmt.Println(&Pair[int, int]{Key: 1, Value: 2})
}

func genericStructInMap() {
	// As map value.
	fmt.Println(map[string]Container[int]{
		"a": {},          // want "missing required fields: Value"
		"b": {Value: 42},
	})

	// As map key.
	fmt.Println(map[Container[int]]string{
		{Value: 1}:  "one",
		{Label: ""}: "bad", // want "missing required fields: Value"
	})
}

func genericStructInSlice() {
	fmt.Println([]Container[int]{
		{},          // want "missing required fields: Value"
		{Value: 42},
	})
}

func unkeyedGenericLiteral() {
	// Unkeyed literals set all fields, so no diagnostic.
	fmt.Println(Container[int]{42, "label"})
	fmt.Println(Pair[string, int]{"key", 42})
}

func returnGenericWithError() (Container[int], error) {
	switch rand.Int() % 3 {
	case 0:
		return Container[int]{}, errors.New("fail") // ok: non-nil error
	case 1:
		return Container[int]{}, nil // want "missing required fields: Value"
	default:
		return Container[int]{Value: 1}, nil
	}
}

func returnGenericPairWithError() (Pair[string, int], error) {
	if rand.Int()%2 == 0 {
		return Pair[string, int]{}, errors.New("fail") // ok: non-nil error
	}
	return Pair[string, int]{}, nil // want "missing required fields: Key, Value"
}

// Nested generic: a generic struct whose type parameter
// is itself a struct with required fields.
type Wrapper[T any] struct { // want Wrapper:"required<Inner>"
	Inner T // required
	Tag   string
}

func nestedGeneric() {
	fmt.Println(Wrapper[Container[int]]{}) // want "missing required fields: Inner"
	fmt.Println(Wrapper[Container[int]]{
		Inner: Container[int]{}, // want "missing required fields: Value"
	})
	fmt.Println(Wrapper[Container[int]]{
		Inner: Container[int]{Value: 1},
	})
}

type ConcreteNode struct { // want ConcreteNode:"required<ID>"
	ID       int // required
	Children []ConcreteNode
}

func concreteNode() {
	fmt.Println(ConcreteNode{})    // want "missing required fields: ID"
	fmt.Println(ConcreteNode{ID: 1})
}
