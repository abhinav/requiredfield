package requiredfield

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// requiredConfig holds configuration for fields
// marked as required via command-line flags.
type requiredConfig struct {
	requiredFields map[typeSpec][]string // "package/path.Type" -> []Field
}

// parseRequiredConfig parses a requiredfield.rc configuration file
// and returns a requiredConfig.
// Each line in the file should be in the format: "required pkg.Type.Field"
// Empty lines and lines starting with "#" are ignored.
func (c *requiredConfig) Parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		err := func() error {
			key, value, _ := strings.Cut(line, " ")
			switch key {
			case "required":
				if err := c.addRequiredField(value); err != nil {
					return fmt.Errorf("add required field: %w", err)
				}

			default:
				return fmt.Errorf("unknown key %q", key)
			}
			return nil
		}()
		if err != nil {
			return fmt.Errorf("%d:%w", lineNum, err)
		}
	}

	return scanner.Err()
}

// RegisterFlags registers command-line flags
// for configuring required fields.
func (c *requiredConfig) RegisterFlags(flag *flag.FlagSet) {
	flag.Func(
		"required",
		"mark field as required (e.g. pkg.Type.Field); can be specified multiple times",
		c.addRequiredField,
	)

	flag.Func(
		"config",
		"load required field specifications from file; suggested only for standalone usage (not via 'go vet')",
		func(path string) error {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()

			return c.Parse(f)
		},
	)
}

// addRequiredField parses and adds a required field specification.
// The spec must be in the format: package/path.TypeName.FieldName
func (c *requiredConfig) addRequiredField(spec string) error {
	if c.requiredFields == nil {
		c.requiredFields = make(map[typeSpec][]string)
	}

	typeSpec, fieldName, err := parseFieldSpec(spec)
	if err != nil {
		return fmt.Errorf(`expected "package/path.Type.Field": %w`, err)
	}

	c.requiredFields[typeSpec] = append(c.requiredFields[typeSpec], fieldName)
	return nil
}

// RequiredFields returns the list of required field names
// for the given package path and type name.
// Returns nil if no fields are configured for this type.
func (c *requiredConfig) RequiredFields(pkgPath, typeName string) []string {
	if c == nil {
		return nil
	}
	return c.requiredFields[typeSpec{
		packagePath: pkgPath,
		typeName:    typeName,
	}]
}

type typeSpec struct {
	packagePath string
	typeName    string
}

func parseFieldSpec(spec string) (typeSpec, string, error) {
	idx := strings.LastIndex(spec, ".")
	if idx == -1 {
		return typeSpec{}, "", errors.New("no field or type specified")
	}

	spec, fieldName := spec[:idx], spec[idx+1:]
	if fieldName == "" {
		return typeSpec{}, "", errors.New("field name is empty")
	}

	idx = strings.LastIndex(spec, ".")
	if idx == -1 {
		return typeSpec{}, "", errors.New("no package or type specified")
	}

	packagePath, typeName := spec[:idx], spec[idx+1:]
	if packagePath == "" {
		return typeSpec{}, "", errors.New("package path is empty")
	}
	if typeName == "" {
		return typeSpec{}, "", errors.New("type name is empty")
	}

	return typeSpec{
		packagePath: packagePath,
		typeName:    typeName,
	}, fieldName, nil
}
