# Configuration

requiredfield supports loading configuration using the `-config` flag.

## File Format

The configuration file format is line-based.
Each line should be in one of the following forms:

```
key1 value1
# comment
key2 value2
```

Key-value pairs are separated by whitespace.
The supported keys are:

- **required**: Marks a field as required.
  The value must be in the format `package/path.TypeName.FieldName` --
  same as the `-required` flag.

<details>
 <summary>Example</summary>

Create a file named `requiredfield.rc` with the following content:

```
# External package fields
required net/http.Request.Method
required net/http.Request.URL

# Internal package fields
required github.com/example/myapp/config.Config.APIKey
required github.com/example/myapp/config.Config.Database
```

Then run requiredfield with the `-config` flag:

```bash
requiredfield -config requiredfield.rc ./...
```

</details>

Fields specified in the configuration file are merged with:

- Fields marked using `// required` comments in source code
- Fields specified via the `-required` flag

## Usage with `go vet`

While the `-config` flag works with `go vet`,
it is recommended for standalone usage only.

If you do use it with `go vet`,
be aware that the path must be absolute,
and the configuration file will be loaded repeatedly
as `go vet` invokes the linter for each compilation unit.

Therefore, for `go vet`,
prefer using the `-required` flag directly
or `// required` comments in source code.
