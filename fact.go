package requiredfield

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// hasRequiredFields is a Fact attached to structs
// listing its required fields.
type hasRequiredFields struct {
	// List is a list of field names
	// in the struct that are marked required.
	List []string
}

var _ analysis.Fact = (*hasRequiredFields)(nil)

func (*hasRequiredFields) AFact() {}

func (f *hasRequiredFields) String() string {
	return "required<" + strings.Join(f.List, ", ") + ">"
}

// isRequiredField is a Fact attached to fields of anonymous structs
// that are marked required.
type isRequiredField struct{}

var _ analysis.Fact = (*isRequiredField)(nil)

func (*isRequiredField) AFact() {}

func (f *isRequiredField) String() string {
	return "required"
}
