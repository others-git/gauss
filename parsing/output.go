package parsing

import (
	"fmt"
	"regexp"
)

// FormattedDiff difference visualized as object
var FormattedDiff KeySlice

func format(input ConsumableDifference) KeyValue {
	var return_value KeyValue

	FormattedDiff = nil
	/*
		for i := range input["Changed"] {
			path_builder(input["Changed"][i]["Path"].([]string))
		}
		for i := range input["Added"] {
			path_builder(input["Added"][i]["Path"].([]string))
		}
		for i := range input["Removed"] {
			path_builder(input["Removed"][i]["Path"].([]string))

		}
	*/

	return return_value
}

func path_builder(path []string) KeyValue {
	var object KeyValue
	FormattedDiff = nil
	r, _ := regexp.Compile("[0-9]+")
	//path_length := len(path)
	for i := range path {
		if ok, _ := regexp.MatchString("{Index:[0-9]+}", path[i]); ok {
			index := r.FindString(path[i])
			fmt.Println(index)
		} else {

		}
	}

	fmt.Println(path)
	fmt.Println(path)
	return object
}

// Format formatting function
func Format(input ConsumableDifference) KeyValue {
	var return_value KeyValue

	return return_value
}
