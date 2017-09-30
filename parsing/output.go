package parsing

import (
	"fmt"
	"regexp"
)

var FormattedDiff Keyslice

func format(input ConsumableDifference) Keyvalue {
	var return_value Keyvalue

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

func pathBuilder(path []string) Keyvalue {
	var object Keyvalue
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

func Format(input ConsumableDifference) Keyvalue {
	var return_value Keyvalue

	return return_value
}
