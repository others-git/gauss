package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"fmt"
)



/*
func Build(input) parsing.Keyvalue {
	var object parsing.Keyvalue

	for i := range input["Changed"] {
		path_builder(input["Changed"][i]["Path"].([]string))
	}
	for i := range input["Added"] {
		path_builder(input["Added"][i]["Path"].([]string))
	}
	for i := range input["Removed"] {
		path_builder(input["Removed"][i]["Path"].([]string))

	}

	return object
}
*/


func patch(patch parsing.ConsumableDifference, original parsing.Keyvalue) parsing.Keyvalue {


	return parsing.Keyvalue{}
}

// Patch: Creates a new object given a 'patch' and 'original'
func Patch(patch parsing.ConsumableDifference, original parsing.Keyvalue) parsing.Keyvalue {
//	var modified parsing.Keyvalue

	o := original
	path := parsing.PathSplit(patch.Added[0].Path)


	// This actually works lol but not really
	for i := range path {

		r := o[path[i]]
		//fmt.Println(r)
		fmt.Println(path[i])
		o = parsing.Remarshal(r)

	}
	fmt.Println(original)
	fmt.Println(o)


	return parsing.Keyvalue{}
}
