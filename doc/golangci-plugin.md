# Use as a golangci-lint plugin

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
> See [How to add a private linter to golangci-lint][1] for details.

[1]: https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint
