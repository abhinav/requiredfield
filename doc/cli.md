---
absorb: true
---

# Command line usage

requiredfield can be used as a standalone tool or with `go vet`.

```bash
# Standalone
requiredfield ./...

# With go vet
go vet -vettool=$(which requiredfield) ./...
```

## Flags

### `-required`

Mark a field as required without modifying its source code.
This is useful for third-party packages
where you cannot add `// required` comments.

The flag accepts a field specification in the format `package/path.Type.Field`.
You can specify the flag multiple times to mark multiple fields as required.

```bash
# Standalone
requiredfield -required package/path.Type.Field ./...

# With go vet
go vet -vettool=$(which requiredfield) \
  -required package/path.Type.Field \
  ./...
```

> [!NOTE]
>
> Fields marked via `-required` are merged
> with fields marked using `// required` comments.

### `-config`

Load required field specifications from a configuration file.
See [Configuration](config.md) for file format details.

> [!NOTE]
>
> This flag is recommended for standalone usage only.
>
> If used with `go vet`, be aware that
> a) the path must be absolute (not relative),
> and b) the configuration file will be loaded repeatedly
> as `go vet` invokes the linter anew for each compilation unit.

```bash
# Standalone (recommended)
requiredfield -config ./requiredfield.rc ./...

# With go vet (path must be absolute)
go vet -vettool=$(which requiredfield) \
  -config /absolute/path/to/requiredfield.rc \
  ./...
```
