package conf

import (
	"os"
	"testing"

	"github.com/opsgenie/oec/test_util"
	"github.com/stretchr/testify/assert"
)

func TestHttpFieldsFilledCorrectly(t *testing.T) {

	confPath, err := test_util.CreateTempTestFile(mockJsonFileContent, ".json")
	assert.Nil(t, err)

	conf, _ := readFileFromLocal(confPath)

	defer os.Remove(confPath)

	assert.Equal(t, conf.ActionMappings["WithHttpAction"].Flags["url"], "https://opsgenie.com")
	assert.Equal(t, conf.ActionMappings["WithHttpAction"].Flags["method"], "PUT")
	assert.Equal(t, conf.ActionMappings["WithHttpAction"].Flags["headers"], "{\"Authentication\":\"Basic JNjDkNsKaMs\"}")
	assert.Equal(t, conf.ActionMappings["WithHttpAction"].Flags["params"], "{\"Key1\":\"Value1\"}")
}
