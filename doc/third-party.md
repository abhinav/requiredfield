# Required fields in third-party code

If you need to enforce that certain fields are always set
for types you don't control,
use the `-required` flag to mark them as required.
This accepts a field specification in the format:

```
package/path.TypeName.FieldName
```

You can specify the flag multiple times.
For example:

```bash
requiredfield \
  -required net/http.Request.Method \
  -required net/http.Request.URL \
  ./...
```

This marks the `Method` and `URL` fields
of the `net/http.Request` type as required.

```go
package http

type Request struct {
    Method string
    URL    *url.URL
    // ... other fields
}
```

Any code that creates a `Request` without setting these fields will be reported:

```go
req := &http.Request{Header: make(http.Header)}
// ERROR: missing required fields: Method, URL
```

> [!NOTE]
>
> Fields marked via `-required` are merged
> with fields marked using `// required` comments.
