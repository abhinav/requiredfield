# Introduction

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
see [Motivation](MOTIVATION.md).
