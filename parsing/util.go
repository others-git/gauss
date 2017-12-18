package parsing

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"
    "golang.org/x/text/unicode/rangetable"
	"runtime/debug"
)

func marshError(input interface{}, stage string, err error) {
	if err != nil {
		fmt.Println(input)
		fmt.Println(stage)
		debug.PrintStack()
		log.Fatal("Remashalling error! ", err)

	}
}


func Remarshal(input interface{}) KeyValue {
	// This is just a nasty type conversions, marshals an interface and then back into our Keyvalue map type
	var back KeyValue
	out, e := json.Marshal(input)
	marshError(input, "Marshal", e)
	e = json.Unmarshal([]byte(out), &back)
	marshError(input, "Unmarshal", e)
	return back
}


func Slicer(input KeyValue) []string {
	// Creates an array of key names given a Keyvalue map
	var r []string
	for key := range input {
		r = append(r, key)
	}
	return r
}

// PathFormatter: Given an array, construct it into a jmespath expression (string with . separator)
func PathFormatter(input []string) string {
	var r string
	for i := range input {
		str := input[i]
		// Escape a . in string name for parsing later
		str = strings.Replace(str, ".", "\\.", -1)
		if i == (len(input) - 1) {
			r = r + str
		} else {
			r = r + str + "."
		}
	}
	return r
}

// IndexOf: Finds index of an object in a given array
func IndexOf(inputList []string, inputKey string) int {
	for i, v := range inputList {
		if v == inputKey {
			return i
		}
	}
	return -1
}


// UnorderedKeyMatch: Returns a bool dependant on all 'keys' in a map matching.
func UnorderedKeyMatch(o KeyValue, m KeyValue) bool {
	istanbool := true
	oSlice := Slicer(o)
	mSlice := Slicer(m)
	for k := range oSlice {
		val := IndexOf(mSlice, oSlice[k])
		if val == -1 {
			istanbool = false
		}
	}

	for k := range mSlice {
		val := IndexOf(oSlice, mSlice[k])
		if val == -1 {
			istanbool = false
		}
	}
	return istanbool
}

// SliceIndex: Adds an 'index' value to the last string in the slice, used for the 'path' to handle arrays.
func SliceIndex(i int, path []string) []string {

	nPath := make([]string, len(path))
	copy(nPath, path)
	iter := len(nPath) - 1
	nPath[iter] = nPath[iter] + "[" + strconv.Itoa(i) + "]"
	return nPath
}

func MatchAny(compare interface{}, compareSlice []interface{}) bool {
	for i := range compareSlice {
		if reflect.DeepEqual(compare, compareSlice[i]) {
			return true
		}
	}
	return false
}

// DoMapArrayKeysMatch: Uses 'UnorderedKeyMatch' to return a bool for two interfaces if they're both maps
func DoMapArrayKeysMatch(o interface{}, m interface{}) bool {
	if reflect.TypeOf(o).Kind() == reflect.Map && reflect.TypeOf(m).Kind() == reflect.Map {
		return UnorderedKeyMatch(Remarshal(o), Remarshal(m))
	}
	return false
}

// PathSplit: Splits up jmespath format path into a slice, will ignore escaped '.' ; opposite of PathFormatter
func PathSplit(input string) []string {

	str := escape(input)
	for i := range str {
		str[i] = strings.Replace(str[i], "\\.", ".", -1)
	}
	return str
}

func escape(input string) []string {
	slashRange := rangetable.New(rune('\\'))
	dotRange := rangetable.New(rune('.'))
	old := rune(0)
	f := func(c rune) bool {
		switch {
		case old == rune('\\'):
			old = rune(0)
			return false
		case old != rune(0):
			return false
		case unicode.In(c, slashRange):
			old = c
			return false
		default:
			return  unicode.In(c, dotRange)

		}
	}
	return strings.FieldsFunc(input, f)

}

// \ = U+005C
// . = U+002E