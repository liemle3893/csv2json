package config

import (
	"github.com/hashicorp/hcl"
	"github.com/liemle3893/csv2json/json"
	"github.com/liemle3893/csv2json/parser"
	"github.com/pkg/errors"
	"strings"
)

// Config Root Config for file
type Config struct {
	RootPath    string      `hcl:"root"`
	OutPath     string      `hcl:"out_directory"`
	Concurrency rune        `hcl:"concurrency"`
	Directories []Directory `hcl:"directory"`
}

// Validate config.
func (c *Config) validate() error {
	// TODO
	if c.Concurrency == 0 {
		c.Concurrency = 10
	}
	return nil
}

// Directory contains info about directory
type Directory struct {
	Path              string             `hcl:",key"`
	Separator         string             `hcl:"separator"`
	Columns           []ColumnDefinition `hcl:"column"`
	AdditionalColumns []ColumnDefinition `hcl:"additional_column"`
	Skip              bool               `hcl:"skip"`            // Skip this directory
	SkipFirstLine     bool               `hcl:"skip_first_line"` // Skip first line
	IncludePatterns   []string           `hcl:"include"`
	ExcludePatterns   []string           `hcl:"exclude"`
}

// ColumnDefinition column definition
type ColumnDefinition struct {
	Name         string                 `hcl:",key"`
	Type         string                 `hcl:"type"` // Only support Float, Int, String, Boolean, StringArray, IntArray, FloatArray, BooleanArray
	Path         string                 `hcl:"path"`
	Separator    string                 `hcl:"separator"` // Only mean in case of Array
	Skip         bool                   `hcl:"skip"`      // Skip this columns
	DefaultValue string                 `hcl:"default"`   // Must set if type was a additional column
	Indices      map[string]interface{} `hcl:"indices"`
}

// ColumnsDefinition is alias for []ColumnDefinition
type ColumnsDefinition struct {
	columns []ColumnDefinition
}

// ParseConfig parse the given HCL string into a Config struct.
func ParseConfig(hclText string) (*Config, error) {
	result := &Config{}

	hclParseTree, err := hcl.Parse(hclText)
	if err != nil {
		return nil, err
	}
	if err := hcl.DecodeObject(&result, hclParseTree); err != nil {
		return nil, err
	}
	if err := result.validate(); err != nil {
		return nil, err
	}
	return result, nil
}

var nilSlice []string = nil

// Parse csv record into json.JsonObject
func (dir *Directory) Parse(record []string) (json.JsonObject, error) {
	data := json.NewJsonObject()
	additionalColumns := ColumnsDefinition{dir.AdditionalColumns}
	if err := additionalColumns.readRecord(data, nilSlice); err != nil {
		return nil, err
	}
	columns := ColumnsDefinition{dir.Columns}
	if err := columns.readRecord(data, record); err != nil {
		return nil, err
	}
	return data, nil
}

// Read single record into json.JsonObject
func (c *ColumnsDefinition) readRecord(root json.JsonObject, record []string) error {
	for ci, column := range c.columns {
		if column.Skip {
			continue
		}
		pieces := strings.Split(column.Path, ".")
		var currentData json.JsonObject = root
		for i, piece := range pieces {
			if (i + 1) < len(pieces) {
				// Handle JsonNode
				if _, ok := currentData.Get(piece); !ok {
					// Make new JsonPath
					currentData.Put(piece, json.NewJsonObject())
				}
				// Data always exists at this point.
				currentData, _ = currentData.GetObject(piece)
			} else {
				// Handle Json Leaf
				if p, err := parser.FindParser(column.Type); err == nil {
					// Get column value
					var columnValue string
					if ci >= len(record) || len(record[ci]) == 0 {
						// If record does not have enough fields or record[ci] is empty
						// Then columnValue will be default value.
						columnValue = column.DefaultValue
					} else {
						columnValue = record[ci]
					}
					// If indexed --> Data
					if p.IsIndexed() {
						if val, ok := column.Indices[columnValue]; ok {
							currentData[piece] = val
						} else {
							currentData[piece], _ = column.Indices[column.DefaultValue]
						}
					} else {
						// Else parse data
						v, _ := p.Parse(columnValue)
						currentData[piece] = v
					}
				} else {
					return errors.Wrap(err, "Parser not found. "+column.Type)
				}
			}
		}
	}
	return nil
}
