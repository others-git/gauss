package parsing

import (
	"encoding/json"
	"log"
)

type Keyvalue map[string]interface{}
type Keyslice map[string][]Keyvalue

type RemovedDifference struct {
	Key   string
	Path  string
	Value interface{}
}

type AddedDifference struct {
	Key   string
	Path  string
	Value interface{}
}

type ChangedDifference struct {
	Key      string
	Path     string
	NewValue interface{}
	OldValue interface{}
}

type ConsumableDifference struct {
	Changed []ChangedDifference `json:",omitempty"`
	Added   []AddedDifference   `json:",omitempty"`
	Removed []RemovedDifference `json:",omitempty"`
}

func marshError(err error) {
	if err != nil {
		log.Fatal(err)
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
	for k := range ListStripper(o) {
		if IndexOf(ListStripper(m), ListStripper(o)[k]) == -1 {
			istanbool = false
		}
	}
	for k := range ListStripper(m) {
		if IndexOf(ListStripper(o), ListStripper(m)[k]) == -1 {
			istanbool = false
		}
	}
	return istanbool
}
