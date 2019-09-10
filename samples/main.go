package main

import (
	c "github.com/liemle3893/csv2json/config"
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
	log.Printf("\n\n\n\n")

	config.Exec()
}
