# Syntax

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
