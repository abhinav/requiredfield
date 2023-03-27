# Overview

To indicate that a field is required,
add a `// required` comment next to it.

```go
type BufWriter struct {
    W      io.Writer     // required
    Buffer *bytes.Buffer
}
```

This indicates that the `W` field is required.

All instantiations of `BufWriter `using the `T{...}` form
will be required to set the `W` field explicitly.

For example:

```go
w := BufWriter{Buffer: b}
// ERROR: missing required fields: W
```
