package parsing

import (
	"encoding/json"
	"github.com/jmespath/go-jmespath"
)

type Keyvalue map[string]interface{}
type Keyslice map[string][]Keyvalue

// TODO: Update Pathspec to support 'index' for slices.
/*
type Pathspec struct{
	Name string
	Index int
}
*/
type Pathspec []string
func (ps *Pathspec) Indexer() {

}


type RemovedDifference struct{
	Key string
	Path Pathspec
	Value interface{}
}
type AddedDifference struct{
	Key string
	Path Pathspec
	Value interface{}
}
type ChangedDifference struct{
	Key string
	Path Pathspec
	NewValue interface{}
	OldValue interface{}
}

type ConsumableDifference struct{
	Changed []ChangedDifference `json:",omitempty"`
	Added []AddedDifference `json:",omitempty"`
	Removed []RemovedDifference `json:",omitempty"`
}

func Remarshal(input interface{}) Keyvalue {
	var back Keyvalue
	out,_ := json.Marshal(input)
	_ = json.Unmarshal([]byte(out), &back)
	return back
}

func ListStripper(input Keyvalue) []string {
	var r []string
	for key := range input {
		r = append(r, key)
	}
	return r
}

func IndexOf(inputList []string, inputKey string) int {
	for i, v := range inputList {
		if v == inputKey {
			return i
		}
	}
	return -1
}


