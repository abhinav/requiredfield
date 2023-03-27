package c

import "testing"

func TestFoo(t *testing.T) {
	tests := []struct {
		give string // required // want give:"required"
		want string
	}{
		{}, // want "missing required fields: give"
		{give: "foo"},
		{give: "foo", want: "bar"},
		{want: "bar"}, // want "missing required fields: give"
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			// ...
		})
	}
}

func TestBar(t *testing.T) {
	// Test in the same file with inverted conditions.

	tests := []struct {
		give string
		want string // required // want want:"required"
	}{
		{},            // want "missing required fields: want"
		{give: "foo"}, // want "missing required fields: want"
		{give: "foo", want: "bar"},
		{want: "bar"},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			// ...
		})
	}
}
