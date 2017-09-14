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
	"github.com/beard1ess/gauss/ui"
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
		output   string
	}{
		{"addKey_o.json", "addKey_m.json", "addKey_d.json", "machine"},
		{"rmKey_o.json", "rmKey_m.json", "rmKey_d.json", "machine"},
		{"modDeepKey_o.json", "modDeepKey_m.json", "modDeepKey_d.json", "machine"},
		{"addModKey_o.json", "addModKey_m.json", "addModKey_d.json", "machine"},
		{"original.json", "encodedOriginal.json", "noDifference.txt", "machine"},
	}

	/*
	 * Test Logic
	 */
	for _, tc := range testCases {

		t.Run(
			fmt.Sprintf(
				"Origin:%s, Modified:%s, Diff:%s, Output:%s",
				tc.origin,
				tc.modified,
				tc.diff,
				tc.output,
			),

			func(t *testing.T) {

				// Read and unmarshal the expected output.
				expectedJson, err := ioutil.ReadFile("./tests/" + tc.diff)
				require.Nil(err, "The test diff should be readable.")

				var expected parsing.ConsumableDifference
				json.Unmarshal(expectedJson, &expected)
				require.Nil(err, "The test data should be unmarshaled without error.")

				// Execute a Diff against the Origin and Modified test files.
				var testBuffer *bytes.Buffer = bytes.NewBuffer(nil)

				ui.Diff(
					"./tests/"+tc.origin,
					"./tests/"+tc.modified,
					tc.output,
					testBuffer,
				)

				// Read and unmarshal the actual output.
				result, err := ioutil.ReadAll(testBuffer)
				require.Nil(err, "The test buffer should be readable")

				var actual parsing.ConsumableDifference
				json.Unmarshal(result, &actual)
				assert.Nil(err, "The result should be unmarshaled without error.")

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
						string(expectedJson),
						string(result),
					),
				)
			},
		)
	}
}
