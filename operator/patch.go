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

//fmt.Println(original)

//fmt.Println(patch.Added)

    //fmt.Println(parsing.PathSplit(patch.Added[0].Path))

	//fmt.Println(original[parsing.PathSplit(patch.Added[0].Path)[0]])

	var o parsing.Keyvalue
	o = original
	path := parsing.PathSplit(patch.Added[0].Path)

	// This actually works lol
	for p := range path {

		r := o[path[p]]
		fmt.Println(r)
		fmt.Println(path[p])
		o = parsing.Remarshal(r)

	}


	return parsing.Keyvalue{}
}
