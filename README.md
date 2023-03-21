# requiredfield

requiredfield is a linter for Go that helps you ensure
that required fields of a struct are filled when initialized.

To indicate that a field is required,
mark it with a `// required` comment next to it.
For example:

```go
type BufWriter struct {
    W      io.Writer     // required
    Buffer *bytes.Buffer
}
```

This indicates that the `W` field is required.
Following this, all instantiations of `BufWriter `using the `T{...}` form
will be required to set the `W` field explicitly.

## FAQ

### Why a comment instead of a struct tag?

requiredfield is intended to be a compile-time analysis,
and does not want to add runtime artifacts to your application.
Tags on struct fields make it into the compiled binary
and are available at runtime via reflection.

Specifying that a field is required with a comment is simple,
and matches documentation conventions already used in a lot of Go code.
