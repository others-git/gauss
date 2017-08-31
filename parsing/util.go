package parsing

import (
	"encoding/json"
)

type Keyvalue map[string]interface{}
type Keyslice map[string][]Keyvalue


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


