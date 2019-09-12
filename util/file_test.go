package util_test

import (
	"testing"

	"github.com/liemle3893/csv2json/util"
	"github.com/stretchr/testify/assert"
)

func TestRemoveFileExtention(t *testing.T) {
	newFileName := util.RemoveFileExtention("/tmp/abc.xyz")
	assert.Equal(t, "/tmp/abc", newFileName, "File name must be /tmp/abc")
}
