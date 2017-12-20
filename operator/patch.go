package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"fmt"
	"github.com/jmespath/go-jmespath"
)



func patch(patch parsing.ConsumableDifference, original parsing.KeyValue) parsing.KeyValue {


	return parsing.KeyValue{}
}

// Patch: Creates a new object given a 'patch' and 'original'
func Patch(patch parsing.ConsumableDifference, original parsing.KeyValue) parsing.KeyValue {

	for _, i := range patch.Added {
		fmt.Println(i.Path)
		res, _ := jmespath.Search(i.Path, original)
		fmt.Println(res)

	}

	return parsing.KeyValue{}
}
