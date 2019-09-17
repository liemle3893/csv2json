package converter

import (
	"bytes"
	"github.com/liemle3893/csv2json/config"
	_ass "github.com/liemle3893/csv2json/testing/assert"
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
)

func TestFileConverter_Convert(t *testing.T) {

	config, _ := config.ParseConfig(configTxt)
	dirConfig := config.Directories[0]
	c := &directoryConverter{DirectoryConfig: dirConfig}

	reader := strings.NewReader("a_string,false,idx1\na_string,false,idx5")
	writer := &bytes.Buffer{}
	row := c.convert0(wrappedReader{reader, "strings.reader"}, writer)

	// JSON response
	expectedJSON := `{"a":"a_string","b":false,"d":"11"}`
	receivedJson := writer.String()
	t.Logf("%+v, %+v", expectedJSON, receivedJson)
	assert.Equal(t, row, uint32(1), "Processed record should be 1")
	t.Run("Compare result", func(t *testing.T) {
		t.Helper()
		_ass.AreEqualJSON(t, receivedJson, expectedJSON, "Invalid JSON")
	})
}

func BenchmarkConverter_Convert(b *testing.B) {
	configuration, _ := config.ParseConfig(configTxt)
	dirConfig := configuration.Directories[0]
	c := &directoryConverter{DirectoryConfig: dirConfig}
	reader := strings.NewReader("a_string,false,idx1\na_string,false,idx5")
	writer := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		c.convert0(wrappedReader{reader, "strings.reader"}, writer)
	}
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
		default = "idx2"
		indices = { "idx1" = "11", "idx2" = "2" }
		excludes = ["idx5"]
	}
}
`
