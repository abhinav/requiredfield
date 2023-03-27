//go:build tools
// +build tools

package tools

import (
	_ "github.com/mgechev/revive"
	_ "go.abhg.dev/stitchmd"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
