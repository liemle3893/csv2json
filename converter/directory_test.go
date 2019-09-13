package converter

import (
	"bytes"
	"github.com/liemle3893/csv2json/config"
	"github.com/magiconair/properties/assert"
	"reflect"
	"strings"
	"testing"
)

func TestFileConverter_Convert(t *testing.T) {

	config, _ := config.ParseConfig(configTxt)
	dirConfig := config.Directories[0]
	c := &directoryConverter{DirectoryConfig: dirConfig}

	reader := strings.NewReader("a_string,false,idx1")
	writer := &bytes.Buffer{}
	row := c.convert0(reader, writer)

	// JSON response
	expectedJSON := "{\"a\":\"a_string\",\"b\":false,\"d\":\"11\"}"
	receivedJson := writer.String()
	t.Logf("%+v, %+v", reflect.TypeOf(expectedJSON), reflect.TypeOf(receivedJson))
	assert.Equal(t, row, uint32(1), "Processed record should be 1")
	t.Run("Compare result", func(t *testing.T) {
		if receivedJson != receivedJson {
			t.Error("Invalid JSON")
		}
	})
}

var configTxt = `
root = "."
out_directory = "./out"

directory "user_action" {
    include = [ ".*" ]
    exclude = [  ]
	column "a" {
		type = "String"
		path = "a"
		default = "a default value"
	}
	column "b" {
		type = "Boolean"
		path = "b"
	}	
	column "d" {
		type = "Indexed"
		path = "d"
		default = "idx1"
		indices = { "idx1" = "11", "idx2" = "2" }
	}		
}
`
