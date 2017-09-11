package parsing

import (
	"encoding/json"
)

type Keyvalue map[string]interface{}
type Keyslice map[string][]Keyvalue

type RemovedDifference struct{
	Key string
	Path string
	Value interface{}
}

type AddedDifference struct{
	Key string
	Path string
	Value interface{}
}

type ChangedDifference struct{
	Key string
	Path string
	NewValue interface{}
	OldValue interface{}
}

type ConsumableDifference struct{
	Changed []ChangedDifference `json:",omitempty"`
	Added []AddedDifference `json:",omitempty"`
	Removed []RemovedDifference `json:",omitempty"`
}

func Remarshal(input interface{}) Keyvalue {
	// This is just a nasty type conversions, marshals an interface and then back into our Keyvalue map type
	var back Keyvalue
	out,_ := json.Marshal(input)
	_ = json.Unmarshal([]byte(out), &back)
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
		if i == (len(input)-1) {
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


