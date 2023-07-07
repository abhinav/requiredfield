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
