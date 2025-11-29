package requiredfield

import (
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseFieldSpec(t *testing.T) {
	tests := []struct {
		name         string
		spec         string
		wantTypeSpec typeSpec
		wantField    string
	}{
		{
			name: "valid spec",
			spec: "github.com/external/pkg.User.ID",
			wantTypeSpec: typeSpec{
				packagePath: "github.com/external/pkg",
				typeName:    "User",
			},
			wantField: "ID",
		},
		{
			name: "valid spec with nested package",
			spec: "example.com/foo/bar/baz.Config.APIKey",
			wantTypeSpec: typeSpec{
				packagePath: "example.com/foo/bar/baz",
				typeName:    "Config",
			},
			wantField: "APIKey",
		},
		{
			name: "simple package",
			spec: "pkg.Type.Field",
			wantTypeSpec: typeSpec{
				packagePath: "pkg",
				typeName:    "Type",
			},
			wantField: "Field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTypeSpec, gotField, err := parseFieldSpec(tt.spec)
			if err != nil {
				t.Errorf("parseFieldSpec() error = %v, want nil", err)
				return
			}

			if gotTypeSpec != tt.wantTypeSpec {
				t.Errorf("parseFieldSpec() typeSpec = %v, want %v", gotTypeSpec, tt.wantTypeSpec)
			}
			if gotField != tt.wantField {
				t.Errorf("parseFieldSpec() field = %q, want %q", gotField, tt.wantField)
			}
		})
	}
}

func TestParseFieldSpec_Errors(t *testing.T) {
	tests := []struct {
		name    string
		spec    string
		wantErr []string
	}{
		{
			name:    "missing field name",
			spec:    "github.com/pkg.User.",
			wantErr: []string{"field name is empty"},
		},
		{
			name:    "missing type name",
			spec:    "github.com/pkg..Field",
			wantErr: []string{"type name is empty"},
		},
		{
			name:    "no dots",
			spec:    "nopackage",
			wantErr: []string{"no field or type specified"},
		},
		{
			name:    "only one dot",
			spec:    "pkg.Type",
			wantErr: []string{"no package or type specified"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseFieldSpec(tt.spec)

			if err == nil {
				t.Errorf("parseFieldSpec() error = nil, want error")
				return
			}

			errMsg := err.Error()
			for _, substring := range tt.wantErr {
				if !strings.Contains(errMsg, substring) {
					t.Errorf("parseFieldSpec() error = %q, want to contain %q", errMsg, substring)
				}
			}
		})
	}
}

func TestRequiredConfig_AddMultiple(t *testing.T) {
	c := new(requiredConfig)

	// Add multiple fields for the same type.
	if err := c.addRequiredField("pkg.User.ID"); err != nil {
		t.Fatalf("add() error = %v", err)
	}
	if err := c.addRequiredField("pkg.User.Name"); err != nil {
		t.Fatalf("add() error = %v", err)
	}
	if err := c.addRequiredField("pkg.User.Email"); err != nil {
		t.Fatalf("add() error = %v", err)
	}

	fields := c.RequiredFields("pkg", "User")
	if len(fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fields))
	}

	expected := map[string]struct{}{"ID": {}, "Name": {}, "Email": {}}
	for _, f := range fields {
		if _, ok := expected[f]; !ok {
			t.Errorf("unexpected field %q", f)
		}
	}
}

func TestRequiredConfig_RegisterFlags(t *testing.T) {
	c := new(requiredConfig)
	fset := flag.NewFlagSet("test", flag.ContinueOnError)

	c.RegisterFlags(fset)

	args := []string{
		"-required", "pkg.User.ID",
		"-required", "pkg.User.Name",
		"-required", "pkg.Config.APIKey",
	}

	if err := fset.Parse(args); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	tests := []struct {
		pkgPath  string
		typeName string
		want     int
	}{
		{"pkg", "User", 2},
		{"pkg", "Config", 1},
		{"pkg", "Unknown", 0},
	}

	for _, tt := range tests {
		fields := c.RequiredFields(tt.pkgPath, tt.typeName)
		if len(fields) != tt.want {
			t.Errorf("RequiredFields(%q, %q) = %d fields, want %d", tt.pkgPath, tt.typeName, len(fields), tt.want)
		}
	}
}

func TestRequiredConfig_RequiredFields(t *testing.T) {
	c := new(requiredConfig)
	err := errors.Join(
		c.addRequiredField("github.com/pkg.User.ID"),
		c.addRequiredField("github.com/pkg.User.Name"),
		c.addRequiredField("github.com/other.Config.Key"),
	)
	if err != nil {
		t.Fatalf("failed to add required fields: %v", err)
	}

	tests := []struct {
		name     string
		pkgPath  string
		typeName string
		want     []string
	}{
		{
			name:     "existing type with fields",
			pkgPath:  "github.com/pkg",
			typeName: "User",
			want:     []string{"ID", "Name"},
		},
		{
			name:     "existing type with one field",
			pkgPath:  "github.com/other",
			typeName: "Config",
			want:     []string{"Key"},
		},
		{
			name:     "non-existent type",
			pkgPath:  "github.com/pkg",
			typeName: "Unknown",
			want:     nil,
		},
		{
			name:     "non-existent package",
			pkgPath:  "unknown",
			typeName: "Type",
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.RequiredFields(tt.pkgPath, tt.typeName)
			if len(got) != len(tt.want) {
				t.Errorf("RequiredFields() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("RequiredFields() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestRequiredConfig_RequiredFields_nil(t *testing.T) {
	var c *requiredConfig
	got := c.RequiredFields("pkg", "Type")
	if got != nil {
		t.Errorf("nil config should return nil, got %v", got)
	}
}

func TestParseRequiredConfig(t *testing.T) {
	type requiredField struct {
		PkgPath string
		Type    string
		Fields  []string
	}

	tests := []struct {
		name string
		give string

		wantRequiredFields []requiredField
	}{
		{
			name: "valid single field",
			give: joinLines("required pkg.User.ID"),
			wantRequiredFields: []requiredField{
				{PkgPath: "pkg", Type: "User", Fields: []string{"ID"}},
			},
		},
		{
			name: "multiple fields same type",
			give: joinLines(
				"required pkg.User.ID",
				"required pkg.User.Name",
				"required pkg.User.Email",
			),
			wantRequiredFields: []requiredField{
				{PkgPath: "pkg", Type: "User", Fields: []string{"ID", "Name", "Email"}},
			},
		},
		{
			name: "multiple types",
			give: joinLines(
				"required github.com/pkg.User.ID",
				"required github.com/other.Config.Key",
			),
			wantRequiredFields: []requiredField{
				{PkgPath: "github.com/pkg", Type: "User", Fields: []string{"ID"}},
				{PkgPath: "github.com/other", Type: "Config", Fields: []string{"Key"}},
			},
		},
		{
			name: "with comments and empty lines",
			give: joinLines(
				"# This is a comment",
				"required pkg.User.ID",
				"",
				"# Another comment",
				"required pkg.User.Name",
			),
			wantRequiredFields: []requiredField{
				{PkgPath: "pkg", Type: "User", Fields: []string{"ID", "Name"}},
			},
		},
		{
			name: "whitespace handling",
			give: joinLines(
				"  required pkg.User.ID",
				"	required pkg.User.Name",
			),
			wantRequiredFields: []requiredField{
				{PkgPath: "pkg", Type: "User", Fields: []string{"ID", "Name"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(requiredConfig)
			if err := c.Parse(strings.NewReader(tt.give)); err != nil {
				t.Errorf("Parse() error = %v, want nil", err)
				return
			}

			for _, want := range tt.wantRequiredFields {
				got := c.RequiredFields(want.PkgPath, want.Type)
				if len(got) != len(want.Fields) {
					t.Errorf("RequiredFields(%q, %q) = %v, want %v", want.PkgPath, want.Type, got, want.Fields)
					continue
				}
				for i := range got {
					if got[i] != want.Fields[i] {
						t.Errorf("RequiredFields(%q, %q) = %v, want %v", want.PkgPath, want.Type, got, want.Fields)
						break
					}
				}
			}
		})
	}
}

func TestParseRequiredConfig_errors(t *testing.T) {
	tests := []struct {
		name    string
		give    string
		wantErr []string
	}{
		{
			name:    "UnknownKey",
			give:    joinLines("optional pkg.User.ID"),
			wantErr: []string{`1:unknown key "optional"`},
		},
		{
			name:    "BadFieldSpec/NoFieldOrType",
			give:    joinLines("required invalid"),
			wantErr: []string{"1:", "no field or type specified"},
		},
		{
			name:    "BadFieldSpec/NoField",
			give:    joinLines("required pkg.User"),
			wantErr: []string{"1:", "no package or type specified"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := new(requiredConfig).Parse(strings.NewReader(tt.give))

			if err == nil {
				t.Errorf("Parse() error = nil, want error")
				return
			}

			errMsg := err.Error()
			for _, substring := range tt.wantErr {
				if !strings.Contains(errMsg, substring) {
					t.Errorf("Parse() error = %q, want to contain %q", errMsg, substring)
				}
			}
		})
	}
}

func TestRequiredConfig_ConfigFlag(t *testing.T) {
	tests := []struct {
		name string
		give string
		want map[typeSpec][]string
	}{
		{
			name: "valid single field",
			give: joinLines("required pkg.User.ID"),
			want: map[typeSpec][]string{
				{packagePath: "pkg", typeName: "User"}: {"ID"},
			},
		},
		{
			name: "multiple fields same type",
			give: joinLines(
				"required pkg.User.ID",
				"required pkg.User.Name",
			),
			want: map[typeSpec][]string{
				{packagePath: "pkg", typeName: "User"}: {"ID", "Name"},
			},
		},
		{
			name: "multiple types",
			give: joinLines(
				"required pkg.User.ID",
				"required pkg.Config.Key",
			),
			want: map[typeSpec][]string{
				{packagePath: "pkg", typeName: "User"}:   {"ID"},
				{packagePath: "pkg", typeName: "Config"}: {"Key"},
			},
		},
		{
			name: "with comments and empty lines",
			give: joinLines(
				"# comment",
				"required pkg.User.ID",
				"",
				"required pkg.Config.Key",
			),
			want: map[typeSpec][]string{
				{packagePath: "pkg", typeName: "User"}:   {"ID"},
				{packagePath: "pkg", typeName: "Config"}: {"Key"},
			},
		},
		{
			name: "empty file",
			want: map[typeSpec][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), "config.rc")

			if err := os.WriteFile(configPath, []byte(tt.give), 0o600); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}

			c := new(requiredConfig)
			fset := flag.NewFlagSet("test", flag.ContinueOnError)
			c.RegisterFlags(fset)

			args := []string{"-config", configPath}
			if err := fset.Parse(args); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if tt.want == nil {
				tt.want = map[typeSpec][]string{}
			}
			if c.requiredFields == nil {
				c.requiredFields = map[typeSpec][]string{}
			}

			if !reflect.DeepEqual(c.requiredFields, tt.want) {
				t.Errorf("requiredFields = %v, want %v", c.requiredFields, tt.want)
			}
		})
	}
}

func TestRequiredConfig_ConfigFlag_errors(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr []string
	}{
		{
			name:    "nonexistent file",
			path:    "/nonexistent/requiredfield-test-file.rc",
			wantErr: []string{"no such file"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(requiredConfig)
			fset := flag.NewFlagSet("test", flag.ContinueOnError)
			fset.SetOutput(io.Discard)
			c.RegisterFlags(fset)

			args := []string{"-config", tt.path}
			err := fset.Parse(args)

			if err == nil {
				t.Errorf("Parse() error = nil, want error")
				return
			}

			errMsg := err.Error()
			for _, substring := range tt.wantErr {
				if !strings.Contains(errMsg, substring) {
					t.Errorf("Parse() error = %q, want to contain %q", errMsg, substring)
				}
			}
		})
	}
}

func TestRequiredConfig_configFlagMergeWithRequired(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.rc")

	content := joinLines(
		"required pkg.User.ID",
		"required pkg.Config.Key",
	)
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	c := new(requiredConfig)
	fset := flag.NewFlagSet("test", flag.ContinueOnError)
	c.RegisterFlags(fset)

	args := []string{
		"-config", configPath,
		"-required", "pkg.User.Name",
		"-required", "other.Type.Field",
	}
	if err := fset.Parse(args); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	want := map[typeSpec][]string{
		{packagePath: "pkg", typeName: "User"}:   {"ID", "Name"},
		{packagePath: "pkg", typeName: "Config"}: {"Key"},
		{packagePath: "other", typeName: "Type"}: {"Field"},
	}

	if !reflect.DeepEqual(c.requiredFields, want) {
		t.Errorf("requiredFields = %v, want %v", c.requiredFields, want)
	}
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}
