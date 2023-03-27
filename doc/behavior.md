# Behavior

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
