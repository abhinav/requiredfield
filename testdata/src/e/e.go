package e

import (
	"errors"
	"fmt"
	"math/rand"
)

type Foo struct { // want Foo:"required<Bar>"
	Bar int // required
	Baz string
}

type MyError struct { // want MyError:"required<Msg>"
	Msg string // required
}

func (e *MyError) Error() string { return e.Msg }

var errSadness = errors.New("great sadness")

func noReturnValue() {
	if rand.Int()%2 == 0 {
		return // ok
	}
	fmt.Println(Foo{}) // want "missing required fields: Bar"
}

func nonNilError() (Foo, error) {
	if rand.Int()%2 == 0 {
		return Foo{}, errors.New("fail") // ok
	} else {
		return Foo{}, errSadness // ok
	}
}

func nilError() (Foo, error) {
	if rand.Int()%2 == 0 {
		return Foo{Bar: 42}, nil // ok
	} else {
		return Foo{}, nil // want "missing required fields: Bar"
	}
}

func errorObject() (Foo, error) {
	if rand.Int()%2 == 0 {
		return Foo{}, &MyError{Msg: "great sadness"} // ok
	} else {
		return Foo{}, &MyError{} // want "missing required fields: Msg"
	}
}

var functionLiteral = func() (Foo, error) {
	if rand.Int()%2 == 0 {
		return Foo{}, &MyError{Msg: "great sadness"} // ok
	} else {
		return Foo{}, &MyError{} // want "missing required fields: Msg"
	}
}

func errorIsNotLastReturn() (error, Foo) {
	if rand.Int()%2 == 0 {
		return errors.New("fail"), Foo{} // want "missing required fields: Bar"
	} else {
		return nil, Foo{} // want "missing required fields: Bar"
	}
}

func errorObjectSubexpression() (Foo, error) {
	switch rand.Int() % 3 {
	case 0:
		return Foo{}, &MyError{Msg: "great sadness"} // ok
	case 1:
		return Foo{}, fmt.Errorf("%w", &MyError{}) // want "missing required fields: Msg"
	default:
		return Foo{}, fmt.Errorf("%w", &MyError{Msg: "great sadness"}) // ok
	}
}
