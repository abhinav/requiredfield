# requiredfield

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
  - [From the command line](#command-line-usage)
  - [As a golangci-lint plugin](#use-as-a-golangci-lint-plugin)
- [Overview](#overview)
  - [Syntax](#syntax)
  - [Behavior](#behavior)
- [FAQ](#faq)
- [Motivation](#motivation)
- [License](#license)

## Introduction

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
see [Motivation](#motivation).

## Installation

Install the binary from source by running:

```bash
go install go.abhg.dev/requiredfield/cmd/requiredfield@latest
```

## Usage

### Command line usage

To use the linter, run the linter with `go vet`:

```bash
go vet -vettool=$(which requiredfield) ./...
```

### Use as a golangci-lint plugin

To use requiredfield as a golangci-lint plugin,
take the following steps:

- Clone the repository or download a source archive
  from the Releases page.

  ```bash
  git clone https://github.com/abhinav/requiredfield.git
  ```

- Build the plugin:

  ```bash
  cd requiredfield
  go build -buildmode=plugin ./cmd/requiredfield
  ```

- Add the linter under `linters-settings.custom` in your `.golangci.yml`,
  referring to the compiled plugin (usually called 'requiredfield.so').

  ```yaml
  linters-settings:
    custom:
      requiredfield:
        path: requiredfield.so
        description: Checks for required struct fields.
        original-url: go.abhg.dev/requiredfield
  ```

- Run golangci-lint as usual.

> **Warning**:
>
> For this to work, your plugin must be built for the same environment
> as the golangci-lint binary you're using.
>
> See [How to add a private linter to golangci-lint](https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint) for details.

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

### Syntax

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

type Post struct {
    Title string // required
}
```

The description is for the benefit of other readers only.
requiredfield will ignore it.

If a field list defines multiple fields on the same line,
all fields will be marked as required.

```go
type City struct {
    Name, State string // required
    Population  int
}
```

#### Positioning

The `// required` comment must be on the line where the field is defined.

```
GOOD                         | BAD
-----------------------------+-------------------
type User struct {           | type User struct {
    Name string // required  |     // required
}                            |     Name string
                             | }
```

If the field definition is spread across multiple lines,
the comment must be on the last of these.
For example,

```go
type Watcher struct {
    Callback func(
        ctx context.Context,
        req *Request,
    ) // required
}
```

Note that you can still place documentation comments for the field above it;
this will not conflict with the `// required` comment.

```go
type User struct {
   // Name is the name of the user.
   Name string // required
}

type Watcher struct {
    // Callback is the function that the Watcher will invoke
    // after it processes a request.
    Callback func(
        ctx context.Context,
        req *Request,
    ) // required
}
```

### Behavior

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

This behavior is not enforced if the struct is being initialized
as part of a return statement with a (probably) non-nil error value:

```go
if err != nil {
    return User{}, err // ok, because the error is non-nil
}
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

### Where does documentation for a field go?

Place documentation for a field above it as you normally would.
This will not affect requiredfield's behavior.

```go
type User struct {
   // Name is the name of the user.
   Name string // required
}
```

## Motivation

A common pattern in Go is to use a struct
to pass several parameters to a function.
This is often referred to as a "parameter object" or a "parameter struct".
If you're unfamiliar with the concept, you can read more about it in
[Designing Go Libraries > Parameter objects](https://abhinavg.net/2022/12/06/designing-go-libraries/#parameter-objects).

In short, the pattern provides some advantages:

1. readability:
   names of fields are visible at call sites,
   allowing them to act as a form of documentation
   similar to named parameters in other languages
2. flexibility:
   new fields can be added without updating all existing call sites

These are both desirable properties for libraries:
users of the library get a readable API
and maintainers of the library can add new **optional** fields
without a major version bump.

For applications, however,
the flexibility afforded by the pattern can turn into a problem.
Application-internal packages rarely cares about API backwards compatibility
and are prone to adding new *required* parameters to functions.
If they use parameter objects,
they lose the ability to safely add these required parameters:
they can no longer have the compiler tell them that they missed a spot.

So application developers are left to choose between:

- parameter objects: get readability, lose safety
- functions with tens of parameters: lose readability, get safety

requiredfield aims to fill this gap with parameter objects
so that applications can still get the readability benefits of using them
without sacrificing safety.

## License

This software is made available under the BSD3 license.
