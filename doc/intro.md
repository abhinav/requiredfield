# Introduction

[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/requiredfield.svg)](https://pkg.go.dev/go.abhg.dev/requiredfield)
[![CI](https://github.com/abhinav/requiredfield/actions/workflows/ci.yml/badge.svg)](https://github.com/abhinav/requiredfield/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/abhinav/requiredfield/branch/main/graph/badge.svg?token=8UW2S4GBTF)](https://codecov.io/github/abhinav/requiredfield)

requiredfield is a linter for Go that verifies that
required fields of a struct are filled when it is initialized.
Whether a field is required is specified with a comment.

For example:

```go
type BufWriter struct {
    W      io.Writer     // required
    Buffer *bytes.Buffer
}
```

The linter will return an error on the following snippet:

```go
w := BufWriter{Buffer: b}
// ERROR: missing required fields: W
```

To read more about the motivation for this linter,
see [Motivation](motivation.md).
