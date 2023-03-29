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

## Positioning

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
