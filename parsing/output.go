package parsing

import (
	"fmt"
	"regexp"
)

// FormattedDiff difference visualized as object
var FormattedDiff KeySlice

func format(input ConsumableDifference) KeyValue {
	var returnValue KeyValue

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

	return returnValue
}


func pathBuilder(input string) KeyValue {
	var object KeyValue
	FormattedDiff = nil
	r, _ := regexp.Compile("[0-9]+")

	path := PathSplit(input)

	//path_length := len(path)
	for i := range path {
		if ok, _ := regexp.MatchString("{Index:[0-9]+}", path[i]); ok {
			index := r.FindString(path[i])
			fmt.Println(index)
		} else {

		}
	}

	fmt.Println(path)
	fmt.Println(len(path))
	return object
}

// Format formatting function
func Format(input ConsumableDifference) KeyValue {
	var returnValue KeyValue

	pathBuilder(input.Added[0].Path)

	return returnValue
}
