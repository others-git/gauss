package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/beard1ess/gauss/parsing"
	"github.com/beard1ess/gauss/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {

	/*
	 * Create instances of test helper objects
	 */
	assert := assert.New(t)
	require := require.New(t)

	/*
	 * Test Cases
	 */
	testCases := []struct {
		origin   string
		modified string
		diff     string
		path     string
		output   string
	}{
		{"addKey_o.json", "addKey_m.json", "addKey_d.json", "", "raw"},
		{"escapeDiff_o.json", "escapeDiff_m.json", "escapeDiff_d.json", "", "raw"},
		{"rmKey_o.json", "rmKey_m.json", "rmKey_d.json", "", "raw"},
		{"modDeepKey_o.json", "modDeepKey_m.json", "modDeepKey_d.json", "", "raw"},
		{"addModKey_o.json", "addModKey_m.json", "addModKey_d.json", "", "raw"},
		{"modArray_o.json", "modArray_m.json", "modArray_d.json", "", "raw"},
		{"modDeepKey_o.json", "modDeepKey_m.json", "modPathDeepKey_d.json", "key1.\"key1-1\"", "raw"},
		{"modArrayKeepPath_o.json", "modArrayKeepPath_m.json", "modArrayKeepPath_d.json", "Outputs.Default[5]", "raw"},
		{"cfn_scaledTest_o.json","cfn_scaledTest_m.json", "cfn_scaledTest_d.json", "", "raw"},
	}
//
	/*
	 * Test Logic
	 */
	for _, tc := range testCases {

		t.Run(
			fmt.Sprintf(
				"Origin:%s, Modified:%s, Diff:%s, Path:%s, Output:%s",
				tc.origin,
				tc.modified,
				tc.diff,
				tc.path,
				tc.output,
			),

			func(t *testing.T) {

				// ReadFile and unmarshal the expected output.
				expectedJSON, err := ioutil.ReadFile("./tests/" + tc.diff)
				require.Nil(err, "The test diff should be readable.")

				var expected parsing.ConsumableDifference
				json.Unmarshal(expectedJSON, &expected)
				require.Nil(err, "The test data should be unmarshaled without error.")
				// sort the arrays

				// Execute a Diff against the Origin and Modified test files.
				var testBuffer *bytes.Buffer = bytes.NewBuffer(nil)

				ui.Diff(
					"./tests/"+tc.origin,
					"./tests/"+tc.modified,
					tc.output,
					tc.path,
					"",
					testBuffer,
				)

				// ReadFile and unmarshal the actual output.
				result, err := ioutil.ReadAll(testBuffer)
				require.Nil(err, "The test buffer should be readable")

				var actual parsing.ConsumableDifference
				json.Unmarshal(result, &actual)
				assert.Nil(err, "The result should be unmarshaled without error.")

				// sort the arrays
				actual.Sort()
				expected.Sort()

				// The Expected and Actual output should be deeply equal.
				assert.Equal(
					reflect.DeepEqual(expected, actual),
					true,
					fmt.Sprintf(
						strings.Join(
							[]string{
								"The diff of %s and %s should equal %s.\n",
								"Expected:\n",
								"%s",
								"Actual:\n",
								"%s",
							},
							" ",
						),
						tc.origin,
						tc.modified,
						tc.diff,
						string(expectedJSON),
						string(result),
					),
				)
			},
		)
	}
}


func TestPatch(t *testing.T) {


	 // Create instances of test helper objects

	assert := assert.New(t)
	require := require.New(t)


	 // Test Cases

	testCases := []struct {
		origin   string
		modified string
		patch     string
		skipKeys     string
		output   string
	}{
		{"addKey_o.json", "addKey_m.json", "addKey_d.json", "", "raw"},
		{"escapeDiff_o.json", "escapeDiff_m.json", "escapeDiff_d.json", "", "raw"},
		{"rmKey_o.json", "rmKey_m.json", "rmKey_d.json", "", "raw"},
		{"modDeepKey_o.json", "modDeepKey_m.json", "modDeepKey_d.json", "", "raw"},
		{"addModKey_o.json", "addModKey_m.json", "addModKey_d.json", "", "raw"},
		{"modArray_o.json", "modArray_m.json", "modArray_d.json", "", "raw"},
		{"modDeepKey_o.json", "modDeepKey_m.json", "modPathDeepKey_d.json", "key1.\"key1-1\"", "raw"},
		{"modArrayKeepPath_o.json", "modArrayKeepPath_m.json", "modArrayKeepPath_d.json", "Outputs.Default[5]", "raw"},
		{"cfn_scaledTest_o.json","cfn_scaledTest_m.json", "cfn_scaledTest_d.json", "", "raw"},
	}
	//

	 // Test Logic

	for _, tc := range testCases {

		t.Run(
			fmt.Sprintf(
				"Origin:%s, Modified:%s, Diff:%s, Path:%s, Output:%s",
				tc.origin,
				tc.modified,
				tc.patch,
				tc.skipKeys,
				tc.output,
			),

			func(t *testing.T) {

				// ReadFile and unmarshal the expected output.
				expectedJSON, err := ioutil.ReadFile("./tests/" + tc.modified)
				require.Nil(err, "The test diff should be readable.")

				var expected parsing.ConsumableDifference
				json.Unmarshal(expectedJSON, &expected)
				require.Nil(err, "The test data should be unmarshaled without error.")
				// sort the arrays

				// Execute a Diff against the Origin and Modified test files.
				var testBuffer = bytes.NewBuffer(nil)

				ui.Patch(
					"./tests/"+tc.modified,
					"./tests/"+tc.origin,
					tc.output,
					tc.skipKeys,
					"",
					"",
					testBuffer,
				)

				// ReadFile and unmarshal the actual output.
				result, err := ioutil.ReadAll(testBuffer)
				require.Nil(err, "The test buffer should be readable")

				var actual parsing.ConsumableDifference
                json.Unmarshal(result, &actual)
				assert.Nil(err, "The result should be unmarshaled without error.")

				// sort the arrays
				actual.Sort()
				expected.Sort()

				// The Expected and Actual output should be deeply equal.
				assert.Equal(
					reflect.DeepEqual(expected, actual),
					true,
					fmt.Sprintf(
						strings.Join(
							[]string{
								"The diff of %s and %s should equal %s.\n",
								"Expected:\n",
								"%s",
								"Actual:\n",
								"%s",
							},
							" ",
						),
						tc.origin,
						tc.modified,
						tc.patch,
						string(expectedJSON),
						string(result),
					),
				)
			},
		)
	}
}
