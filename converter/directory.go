package converter

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"github.com/liemle3893/csv2json/config"
	"github.com/liemle3893/csv2json/util"
	"io"
	"log"
	"os"
	"path"
)

type directoryConverter struct {
	CSVDirectory    string
	JSONDirectory   string
	DirectoryConfig config.Directory
}

func newDirectoryConverter(csvDirectory, jsonDirectory string, directoryConfig config.Directory) *directoryConverter {
	return &directoryConverter{csvDirectory, jsonDirectory, directoryConfig}
}

func (c *directoryConverter) convert(fileCountChannel, rowCountChannel chan uint32) error {
	dirConfig := c.DirectoryConfig
	if dirConfig.Skip {
		return nil
	}
	csvDirName := c.CSVDirectory
	jsonDirName := c.JSONDirectory
	err := os.MkdirAll(jsonDirName, os.ModePerm)
	if err != nil {
		return err
	}
	files, _ := util.ListFiles(csvDirName, dirConfig.IncludePatterns, dirConfig.ExcludePatterns)
	for _, file := range files {
		rows, err := c.convertFile(csvDirName, jsonDirName, file)
		if err != nil {
			log.Printf("%+v", err)
		}
		fileCountChannel <- uint32(1)
		rowCountChannel <- rows
	}
	return nil
}

func (c *directoryConverter) convertFile(csvDirName, jsonDirName string, file os.FileInfo) (uint32, error) {
	fileName := file.Name()
	csvFilePath := path.Join(csvDirName, fileName)
	jsonFilePath := util.RemoveFileExtention(path.Join(jsonDirName, fileName)) + ".json"
	csvReader, err := os.Open(csvFilePath)
	defer csvReader.Close()
	if err != nil {
		log.Printf("Fail to open file: %s. %+v\n", csvFilePath, err)
		return 0, err
	}
	jsonWriter, err := os.Create(jsonFilePath)
	defer jsonWriter.Close()
	if err != nil {
		log.Printf("Fail to create file: %s. %+v\n", jsonFilePath, err)
		return 0, err
	}
	w := bufio.NewWriterSize(jsonWriter, 4096*2) // 8KB buf size
	rows := c.convert0(bufio.NewReader(csvReader), w)
	if err := w.Flush(); err != nil {
		return rows, err
	}
	return rows, err
}

func (c *directoryConverter) createCsvReader(reader io.Reader) *csv.Reader {
	csvReader := csv.NewReader(reader)
	dirConfig := c.DirectoryConfig
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = -1 // Variables fields
	if "" != dirConfig.Separator {
		csvReader.Comma = rune(dirConfig.Separator[0])
	}
	return csvReader
}

func (c *directoryConverter) convert0(csvReader io.Reader, jsonWriter io.StringWriter) uint32 {
	var rowCount uint32 = 0
	var firstLine = true
	r := c.createCsvReader(csvReader)
	for {
		if firstLine && c.DirectoryConfig.SkipFirstLine {
			firstLine = false
			continue
		}
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		data, err := c.DirectoryConfig.Parse(record)
		if err != nil {
			log.Println(err)
			continue
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = jsonWriter.WriteString(string(jsonData) + "\n")
		if err != nil {
			log.Printf("Fail to write file. %+v\n", err)
		}
		rowCount++
	}
	return rowCount
}
