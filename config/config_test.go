package config_test

import (
	"reflect"
	"testing"

	"github.com/liemle3893/csv2json/config"
)

func TestConfigParsing(t *testing.T) {
	expected := &config.Config{
		RootPath: ".",
		OutPath:  "./out",
		Directories: []config.Directory{
			config.Directory{
				Path:            "user_action",
				Separator:       "",
				IncludePatterns: []string{".*"},
				ExcludePatterns: []string{},
				Columns: []config.ColumnDefinition{
					config.ColumnDefinition{Name: "a", Type: "String", DefaultValue: "a default value", Path: "a"},
					config.ColumnDefinition{Name: "b", Type: "Boolean", Skip: true, Path: "b"},
					config.ColumnDefinition{Name: "d", Type: "Indexed", Skip: true, Path: "b", Indices: map[string]interface{}{
						"idx1": "1",
						"idx2": "2",
					}, DefaultValue: "idx1"},
				},
			},
		},
	}

	config, err := config.ParseConfig(configTxt)

	t.Logf("%+v\n", config)
	t.Logf("%+v\n", expected)
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
	column "d" {
		type = "Indexed"
		path = "b"
		skip = true
		default = "idx1"
		indices = { "idx1" = "1", "idx2" = "2" }
	}		
}
`
