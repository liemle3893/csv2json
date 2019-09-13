package config

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
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
	var result error = nil
	// TODO
	if c.Concurrency == 0 {
		c.Concurrency = 10
	}
	for idx, _ := range c.Directories {
		dir := &(c.Directories)[idx]
		err := dir.validate()
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	fmt.Printf("Result: %+v", result)
	return result
}

// Directory contains info about directory
type Directory struct {
	Path              string              `hcl:",key"`
	Separator         string              `hcl:"separator"`
	Columns           []*ColumnDefinition `hcl:"column"`
	AdditionalColumns []*ColumnDefinition `hcl:"additional_column"`
	Skip              bool                `hcl:"skip"`            // Skip this directory
	SkipFirstLine     bool                `hcl:"skip_first_line"` // Skip first line
	IncludePatterns   []string            `hcl:"include"`
	ExcludePatterns   []string            `hcl:"exclude"`
}

func (dir *Directory) validate() error {
	var result error = nil
	for _, column := range dir.Columns {
		err := column.validate()
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
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
	Excludes     []string               `hcl:"excludes"` // Exclude value
	excludeMap   map[string]bool
}

func (c *ColumnDefinition) validate() error {
	if len(c.Excludes) > 0 {
		c.excludeMap = make(map[string]bool)
		for _, exclude := range c.Excludes {
			c.excludeMap[exclude] = true
		}
	}
	return nil
}
func (c *ColumnDefinition) shouldSkip(value string) bool {
	_, ok := c.excludeMap[value]
	return ok
}

// ColumnsDefinition is alias for []ColumnDefinition
type ColumnsDefinition struct {
	columns []*ColumnDefinition
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
	if _, err := additionalColumns.readRecord(data, nilSlice); err != nil {
		return nil, err
	}
	columns := ColumnsDefinition{dir.Columns}
	ok, err := columns.readRecord(data, record)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return data, nil
}

// Read single record into json.JsonObject.
// Return true if record read success. False to skip read.
func (c *ColumnsDefinition) readRecord(root json.JsonObject, record []string) (bool, error) {
	for ci, column := range c.columns {
		if column.Skip {
			continue
		}
		pieces := strings.Split(column.Path, ".")
		var currentData = root
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
					if column.shouldSkip(columnValue) {
						return false, nil
					}
					// If indexed --> Data
					if p.IsIndexed() {
						if val, ok := column.Indices[columnValue]; ok {
							currentData.Put(piece, val)
						} else {
							currentData.Put(piece, column.Indices[column.DefaultValue])
						}
					} else {
						// Else parse data
						v, _ := p.Parse(columnValue)
						currentData.Put(piece, v)
					}
				} else {
					return true, errors.Wrap(err, "Parser not found. "+column.Type)
				}
			}
		}
	}
	return true, nil
}
