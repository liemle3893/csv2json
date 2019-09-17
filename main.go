package main

import (
	"github.com/liemle3893/csv2json/cmd"
	"log"
)


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
