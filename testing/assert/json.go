package assert

import (
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func AreEqualJSON(t *testing.T, got, want string, msg ...string) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(got), &o1)
	if err != nil {
		fmt.Printf("Error mashalling string 1 :: %s", err.Error())
		t.Fail()
	}
	err = json.Unmarshal([]byte(want), &o2)
	if err != nil {
		fmt.Printf("Error mashalling string 2 :: %s", err.Error())
		t.Fail()
	}
	assert.Equal(t, o1, o2, msg...)
}
