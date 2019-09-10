package config

import (
	"reflect"
	"testing"
)

func TestConfigParsing(t *testing.T) {
	expected := &Config{
		RootPath: ".",
		OutPath:  "./out",
		Directories: []Directory{
			Directory{
				Path:            "user_action",
				Separator:       "",
				IncludePatterns: []string{".*"},
				ExcludePatterns: []string{},
				Columns: []ColumnDefinition{
					ColumnDefinition{Name: "a", Type: "String", DefaultValue: "a default value", Path: "a"},
					ColumnDefinition{Name: "b", Type: "Boolean", Skip: true, Path: "b"},
				},
			},
		},
	}

	config, err := ParseConfig(configTxt)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(config, expected) {
		t.Error("Config structure differed from expectation")
	}
}

var configTxt = `
root = "."
out_directory = "./out"

directory "user_action" {
    include = [ ".*" ]
    exclude = [  ]
	column "a" {
		type = "String"
		path = "a"
		default = "a default value"
	}
	column "b" {
		type = "Boolean"
		path = "b"
		skip = true
	}	
}
`
