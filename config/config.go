package config

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liemle3893/csv2json/parser"
	"github.com/liemle3893/csv2json/util"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/pkg/errors"
)

// Config Root Config for file
type Config struct {
	RootPath    string      `hcl:"root"`
	OutPath     string      `hcl:"out_directory"`
	Directories []Directory `hcl:"directory"`
}

// Directory contains info about directory
type Directory struct {
	Path              string             `hcl:",key"`
	Separator         string             `hcl:"separator"`
	Columns           []ColumnDefinition `hcl:"column"`
	AdditionalColumns []ColumnDefinition `hcl:"additional_column"`
	Skip              bool               `hcl:"skip"`        // Skip this directory
	SkipHeader        bool               `hcl:"skip_header"` // Skip first line
	IncludePatterns   []string           `hcl:"include"`
	ExcludePatterns   []string           `hcl:"exclude"`
}

// ColumnDefinition column definition
type ColumnDefinition struct {
	Name         string `hcl:",key"`
	Type         string `hcl:"type"` // Only support Float, Int, String, Boolean, StringArray, IntArray, FloatArray, BooleanArray
	Path         string `hcl:"path"`
	Separator    string `hcl:"separator"` // Only mean in case of Array
	Skip         bool   `hcl:"skip"`      // Skip this columns
	DefaultValue string `hcl:"default"`   // Must set if type was a additional column
}

// ColumnsDefinition is alias for []ColumnDefinition
type ColumnsDefinition struct {
	columns []ColumnDefinition
}

// ParseConfig parse the given HCL string into a Config struct.
func ParseConfig(hclText string) (*Config, error) {
	result := &Config{}
	var errors *multierror.Error

	hclParseTree, err := hcl.Parse(hclText)
	if err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&result, hclParseTree); err != nil {
		return nil, err
	}
	return result, errors.ErrorOrNil()
}

var nilSlice []string = nil

// Parse csv record into map[string]interface{}
func (dir *Directory) Parse(record []string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
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

func (c *ColumnsDefinition) readRecord(root map[string]interface{}, record []string) error {
	for ci, column := range c.columns {
		if column.Skip {
			continue
		}
		pieces := strings.Split(column.Path, ".")
		var currentData map[string]interface{} = root
		for i, piece := range pieces {
			if (i + 1) < len(pieces) {
				if _, ok := currentData[piece]; !ok {
					currentData[piece] = make(map[string]interface{})
				}
				currentData = currentData[piece].(map[string]interface{})
			} else {
				if p, err := parser.FindParser(column.Type); err == nil {
					var columnValue string
					if ci >= len(record) {
						columnValue = column.DefaultValue
					} else {
						columnValue = record[ci]
					}
					v, _ := p.Parse(columnValue)
					currentData[piece] = v
				} else {
					return errors.Wrap(err, "Parser not found. "+column.Type)
				}
			}
		}
	}
	return nil
}

// Exec export data from CSV to JSON
func (c *Config) Exec() {
	var wg sync.WaitGroup
	var reporter = make(chan int)
	for _, dir := range c.Directories {
		directory := path.Join(c.RootPath, dir.Path)
		outDirectory := path.Join(c.OutPath, dir.Path)
		if err := os.MkdirAll(outDirectory, 0755); err != nil {
			log.Fatalf("Cannot create out directory. %+v", err)
		}
		files, _ := util.ListFiles(directory, dir.IncludePatterns, dir.ExcludePatterns)
		for _, file := range files {
			inFile := path.Join(directory, file.Name())
			var extension = filepath.Ext(file.Name())
			var outFile = path.Join(outDirectory, file.Name())
			outFile = outFile[0:len(outFile)-len(extension)] + ".json"
			wg.Add(1)
			go func(f os.FileInfo) {
				defer wg.Done()
				parseFile(inFile, outFile, dir, f)
				reporter <- 1
			}(file)
		}
	}

	go func() {
		var fileCounter = 0
		for i := range reporter {
			fileCounter += i
			printProcess(fileCounter)
		}
	}()
	wg.Wait()
	close(reporter)
}

func printProcess(count int) {
	// fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\r%d file(s) processed", (count))
}

func parseFile(inputFile, outoutFile string, dir Directory, file os.FileInfo) {
	f, err := os.Open(inputFile)
	if err != nil {
		return
	}
	r := csv.NewReader(bufio.NewReader(f))
	if len(dir.Separator) > 0 {
		r.Comma = rune(dir.Separator[0])
	}
	r.Comment = '#'
	writer, err := os.Create(outoutFile)
	defer writer.Close()
	var firstLine = true
	for {
		if firstLine && dir.SkipHeader {
			firstLine = false
			continue
		}
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		data, err := dir.Parse(record)
		if err != nil {
			log.Println(err)
			continue
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			continue
		}
		writer.WriteString(string(jsonData) + "\n")
	}
}
