package conf

import (
	"github.com/opsgenie/oec/test_util"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadFileFromLocal(t *testing.T) {

	confPath, err := test_util.CreateTempTestFile(mockJsonFileContent, ".json")
	assert.Nil(t, err)

	actualConf, _ := readFileFromLocal(confPath)

	defer os.Remove(confPath)

	assert.Equal(t, mockConf, actualConf,
		"Actual configuration was not equal to expected configuration.")
}
