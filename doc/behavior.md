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

This behavior is not enforced if the struct is being initialized
as part of a return statement with a (probably) non-nil error value:

```go
if err != nil {
    return User{}, err // ok, because the error is non-nil
}
```
