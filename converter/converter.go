package converter

import (
	"fmt"
	"github.com/liemle3893/csv2json/config"
	"log"
	"path"
	"sync"
)

// Converter will read configuration and convert CSV files into JSON files
type Converter struct {
	config *config.Config
}

func NewConverter(config *config.Config) *Converter {
	return &Converter{config: config}
}

// Convert files
func (c *Converter) Convert() {
	dirs := c.config.Directories
	convertersChan := make(chan *directoryConverter, c.config.Concurrency)
	fCountChan, rCountChan := make(chan uint32), make(chan uint32)
	// Concurrently convert file
	var wg sync.WaitGroup
	go func() {
		for dirConverter := range convertersChan {
			err := dirConverter.convert(fCountChan, rCountChan)
			if err != nil {
				log.Printf("Fail to convert directory. %+v", err)
			}
			wg.Done()
		}
	}()
	// Concurrently submit dirConverter
	for _, dir := range dirs {
		csvDir := path.Join(c.config.RootPath, dir.Path)
		jsonDir := path.Join(c.config.OutPath, dir.Path)
		wg.Add(1)
		convertersChan <- newDirectoryConverter(csvDir, jsonDir, dir)
	}
	printProcess(fCountChan, rCountChan)
	// Wait util done.
	wg.Wait()
	close(convertersChan)
	close(fCountChan)
	close(rCountChan)
}

func printProcess(fileCountChannel, rowsCountChannel chan uint32) {
	var fileCount, rowCount = uint32(0), uint32(0)
	_print := func(fCount, rCount uint32) {
		fmt.Printf("\r%d file(s) processed. %d row(s) processed", fCount, rCount)
	}
	go func() {
		for {
			select {
			case f := <-fileCountChannel:
				fileCount += f
				_print(fileCount, rowCount)
			case r := <-rowsCountChannel:
				rowCount += r
				_print(fileCount, rowCount)
			}
		}
	}()
	fileCountChannel <- 0
	rowsCountChannel <- 0
}
