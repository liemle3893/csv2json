package main

import (
	c "encoding/config"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	testConfig, _ := ioutil.ReadFile("config.hcl")
	config, err := c.ParseConfig(string(testConfig))
	if err != nil {
		t.Error(err)
	} else {
		log.Printf("%+v\n", config)
	}
	checkFileContent(t, "out/user_info/test.json")
	checkFileContent(t, "out/user_action/test.json")
}

func checkFileContent(t *testing.T, file string) {
	f, err := os.Open(file)
	if err != nil {
		t.Error(err)
	}
	lines, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	} else {
		if string(lines) != json {
			t.Failed()
		}
	}
}

var json = `{"a":{"s":{"a":64,"c":"1.0","d":true,"ip":"127.0.0.1"},"type":"PING"}}`
