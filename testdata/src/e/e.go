package e

import "fmt"

type Foo struct { // want Foo:"required<A, B>"
	A string // required
	B string // required
	C string
}

func check(f Foo) {
	fmt.Println(f)
}

func _1() {
	check(Foo{}) // want "missing required fields: A, B"

	check(Foo{A: "a"}) // want "missing required fields: B"

	check(Foo{B: "b"}) // want "missing required fields: A"

	check(Foo{A: "a", B: "b"})
}

func _2() {
	var f Foo
	check(f) // want "missing required fields: A, B"
	f.A = "a"
	check(f) // want "missing required fields: B"
	f.B = "b"
	check(f)
}

func _3(cond bool) {
	var f Foo
	if cond {
		f.A = "a"
	}
	f.B = "b"
	check(f) // want "missing required fields: A"
}

func _4(cond bool) {
	var f Foo
	f.A = "a"
	if cond {
		check(f) // want "missing required fields: B"
	}
	f.B = "b"
	check(f)
}

func _5(cond bool) {
	var f Foo
	if cond {
		f.A = "a1"
	} else {
		f.A = "a2"
	}
	f.B = "b"
	check(f)
}

func _6(cond bool) {
	var f Foo
	if cond {
		f.A = "a"
		f.B = "b"
	}
	check(f) // want "missing required fields: A, B"
}

func _7(ok func(int) bool) {
	var f Foo
	for i := 0; i < 10; i++ {
		if ok(i) {
			f.A = "a"
			f.B = "b"
			break
		}
	}
	check(f) // want "missing required fields: A, B"
}
