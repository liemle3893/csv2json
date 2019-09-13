package json

import (
	"github.com/liemle3893/csv2json/testing/assert"
	"testing"
)

func TestJsonObject(t *testing.T) {
	json := NewJsonObject()
	json.Put("Key", "Value")
	json.Put("Int", 1)
	child := NewJsonObject()
	json.Put("child", child)
	c, _ := json.GetObject("child")
	c.Put("Double", 1.5)
	expected := `{"Key":"Value","Int":1,"child":{"Double":1.5}}`
	got := json.String()
	t.Logf("expected: %+v", expected)
	t.Logf("got: %+v", got)
	assert.AreEqualJSON(t, got, expected, "Invalid result")

}
