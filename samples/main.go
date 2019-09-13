package main

import (
	c "github.com/liemle3893/csv2json/config"
	"github.com/liemle3893/csv2json/converter"
	"io/ioutil"
	"log"
)

func main() {
	testConfig, _ := ioutil.ReadFile("config.hcl")
	config, err := c.ParseConfig(string(testConfig))
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("%+v\n", config)
	}
	converter := converter.NewConverter(config)
	converter.Convert()
}
