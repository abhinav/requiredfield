# requiredfield

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
see [Motivation](doc/MOTIVATION.md).

## Installation

Install the binary from source by running:

```bash
go install go.abhg.dev/requiredfield/cmd/requiredfield@latest
```

## Usage

To use the linter, run the binary directly:

```bash
requiredfield ./...
```

Alternatively, use it with `go vet`:

```bash
go vet -vettool=$(which requiredfield) ./...
```

## Overview

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

## Syntax

Fields are marked as required by adding a comment
in one of the following forms next to them:

```go
// required
// required<sep><description>
```

Where `<sep>` is a non-alphanumeric character,
and `<description>` is an optional description.

For example:

```go
type User struct {
    Name  string // required: must be non-empty
    Email string
}
```

The description is for the benefit of other readers only.
requiredfield will ignore it.

## Behavior

Any time a struct is initialized in the form `T{..}`,
requiredfield will ensure that all its required fields are set explicitly.

```go
u := User{
    Email: email,
}
// ERROR: missing required fields: Name
```

Required fields can be set to the zero value of their type,
but that choice must be made explicitly.

```go
u := User{
    Name: "", // computed below
    Email: email,
}
// ...
u.Name = name
```

## FAQ

### Why a comment instead of a struct tag?

The reasons for this choice are both, philosophical and cosmetic.

First, the philosophical reason:
requiredfield is a linter that runs at compile-time,
and therefore wants its footprint limited to compile-time only.
Struct tags get compiled into your binary
and are available at runtime via reflection.
It would become possible for someone to
change how the program behaves based on the value of those struct tags.
requiredfield considers that a violation of the linter's boundaries,
and aims to prevent that by using comments instead.

The cosmetic reason is much easier to explain:
Struct tags are uglier than line comments.

```go
Author ID `required:"true"`

// versus

Author ID // required
```

## License

This software is made available under the MIT license.
