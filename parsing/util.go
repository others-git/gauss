package parsing

import (
	"encoding/json"
	"log"
)



func marshError(err error) {
	if err != nil {
		log.Fatal("Remashalling error! ", err)
	}
}

func Remarshal(input interface{}) Keyvalue {
	// This is just a nasty type conversions, marshals an interface and then back into our Keyvalue map type
	var back Keyvalue
	out, e := json.Marshal(input)
	marshError(e)
	e = json.Unmarshal([]byte(out), &back)
	marshError(e)
	return back
}

func ListStripper(input Keyvalue) []string {
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
	o_slice := ListStripper(o)
	m_slice := ListStripper(m)
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
