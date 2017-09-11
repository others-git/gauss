package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/beard1ess/gauss/parsing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {

	var expected, actual parsing.ConsumableDifference

	assert := assert.New(t)
	require := require.New(t)

	testBuffer := bytes.NewBuffer(nil) // A testing writer

	diff(
		"./tests/one.json",
		"./tests/two.json",
		"machine",
		testBuffer,
	)

	result, err := ioutil.ReadAll(testBuffer)
	require.Nil(err, "The test buffer should be readable")

	testData, err := ioutil.ReadFile("./tests/diff.json")
	require.Nil(err, "The test diff should be readable.")

	json.Unmarshal(testData, &expected)
	require.Nil(err, "The test data should be unmarshaled without error.")

	json.Unmarshal(result, &actual)
	assert.Nil(err, "The result should be unmarshaled without error.")

	assert.Equal(
		reflect.DeepEqual(expected, actual),
		true,
		fmt.Sprintf(
			strings.Join(
				[]string{
					"The diff of one.json and two.json should equal the test diff.\n",
					"Expected:\n",
					"%s",
					"Actual:\n",
					"%s",
				},
				" ",
			),
			string(testData),
			string(result),
		),
	)
}
