package parsing

import (
	"encoding/json"
	"log"
	"fmt"
	"strconv"
	"reflect"
)


func marshError(input interface{}, stage string, err error) {
	if err != nil {
		fmt.Println(input)
		fmt.Println(stage)
		log.Fatal("Remashalling error! ", err)
	}
}

func Remarshal(input interface{}) Keyvalue {
	// This is just a nasty type conversions, marshals an interface and then back into our Keyvalue map type
	var back Keyvalue
	out, e := json.Marshal(input)
	marshError(input, "Marshal", e)
	e = json.Unmarshal([]byte(out), &back)
	marshError(input, "Unmarshal", e)
	return back
}

func Slicer(input Keyvalue) []string {
	// Creates an array of key names given a Keyvalue map
	var r []string
	for key := range input {
		r = append(r, key)
	}
	return r
}

func PathFormatter(input []string) string {
	// Given an array, construct it into a jmespath expression (string with . separator)
	var r string
	for i := range input {
		if i == (len(input) - 1) {
			r = r + input[i]
		} else {
			r = r + input[i] + "."
		}
	}
	return r
}

func IndexOf(inputList []string, inputKey string) int {
	// Finds index of an object given an array
	for i, v := range inputList {
		if v == inputKey {
			return i
		}
	}
	return -1
}

func UnorderedKeyMatch(o Keyvalue, m Keyvalue) bool {
	istanbool := true
	o_slice := Slicer(o)
	m_slice := Slicer(m)
	for k := range o_slice {
		val := IndexOf(m_slice, o_slice[k])
		if val == -1 {
			istanbool = false
		}
	}

	for k := range m_slice {
		val := IndexOf(o_slice, m_slice[k])
		if val == -1 {
			istanbool = false
		}
	}
	return istanbool
}

func PathSlice(i int, path []string ) []string {

	npath := make([]string, len(path))
	copy(npath, path)
	iter := len(npath) - 1
	npath[iter] = npath[iter] + "[" + strconv.Itoa(i) + "]"
	return npath
}

func MatchAny(compare interface{}, compareSlice []interface{}) bool {
	for i := range compareSlice {
		if reflect.DeepEqual(compare, compareSlice[i]) {
			return true
		}
	}
	return false
}