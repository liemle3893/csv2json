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
	Output            string              `hcl:"output"` // Default is Path
	columns           *ColumnsDefinition
	additionalColumns *ColumnsDefinition
}

func (dir *Directory) validate() error {
	var result error = nil
	skipMap := make(map[int]map[string]bool)
	for idx, column := range dir.Columns {
		err := column.validate()
		if err != nil {
			result = multierror.Append(result, err)
		}
		if len(column.Excludes) > 0 {
			excludeMap := make(map[string]bool)
			for _, exclude := range column.Excludes {
				excludeMap[exclude] = true
			}
			skipMap[idx] = excludeMap
		}
	}
	dir.columns = &ColumnsDefinition{dir.Columns, skipMap}
	dir.additionalColumns = &ColumnsDefinition{dir.AdditionalColumns, nil}
	if len(strings.TrimSpace(dir.Output)) == 0 {
		dir.Output = dir.Path
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
}

func (c *ColumnDefinition) validate() error {
	return nil
}

// ColumnsDefinition is alias for []ColumnDefinition
type ColumnsDefinition struct {
	columns []*ColumnDefinition
	skipMap map[int]map[string]bool
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
	additionalColumns := dir.additionalColumns
	if _, err := additionalColumns.readRecord(data, nilSlice); err != nil {
		return nil, err
	}
	columns := dir.columns
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
	if c.skipMap != nil {
		for ci, field := range record {
			excludeMap, ok := c.skipMap[ci]
			if ok {
				if _, ok := excludeMap[field]; ok {
					// There is field that should skip in record
					return false, nil
				}
			}
		}
	}
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
					return true, errors.Wrap(err, "Parser not found for column: "+column.Name+"")
				}
			}
		}
	}
	return true, nil
}
