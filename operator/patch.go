package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"fmt"
	"github.com/jmespath/go-jmespath"
	"encoding/json"
	"regexp"
	"strconv"
	"reflect"
)

// https://github.com/golang/go/wiki/SliceTricks

func patch(patch parsing.ConsumableDifference, original parsing.KeyValue) parsing.KeyValue {


	return parsing.KeyValue{}
}

// Patch: Creates a new object given a 'patch' and 'original'
func Patch(patch *parsing.ConsumableDifference, original *parsing.Gaussian) (*interface{}, error) {
	originalObject := original.Data

	var newObject interface{}
	// Updated order Index > Changed > Removed > Added


	// Iterate over added objects
	for _, i := range patch.Added {

		originPath := i.Path
		key := i.Key
		value := i.Value

		// validate jmespath
		_, err :=  jmespath.Compile(originPath)
		if err != nil {
			nErr := fmt.Errorf("failed to compile provided path: %T", err)
			return nil, nErr
		}

		// slice up path
		slicedPath := parsing.PathSplit(originPath)

		// create child object
		childObject, err := createChild(slicedPath, key, value, originalObject)
		if err != nil {
			return nil, err
		}

		fmt.Println(*childObject)

		// wrap child object to create new object
		newObject, err = addParent(slicedPath, childObject, originalObject)
		if err != nil {
			fmt.Println(err)
		}

		// testing so break
		break
	}

	return &newObject, nil
}

// creates new child object from key and value
func createChild(path []string, key string, value interface{}, object interface{}) (*interface{}, error) {

	var newObject interface{}

	// check path for index
	index, stringPath, err := makePath(path)
	if err != nil {
		return nil, err
	}

	// get working directory based on path
	objectDir,err := jmespath.Search(*stringPath, object)
	if err != nil {
		fmt.Println(*stringPath)
		return nil, err
	}

	// determine what type of object we need to make
	if key != "" {
		newObject = map[string]interface{}{
			key: value,
		}
	} else {
		newObject = value
	}

	// Update logic for slice value
	if index != nil && reflect.TypeOf(objectDir).Kind() == reflect.Slice {
		//TODO: do thing with index
		// cast to slice of interfaces
		objectSlice := objectDir.([]interface{})

		// insert into slice
		objectSlice = append(objectSlice, 0)
		copy(objectSlice[*index+1:], objectSlice[*index:])
		objectSlice[*index] = newObject
		newObject = objectSlice
	}



	return &newObject, nil
}

// recreate the template with the updated child object
func addParent(path []string, child interface{}, object interface{}) (*interface{}, error) {
	// kv type we'll unmarshal to

	if len(path) > 0 {
		var kv map[string]interface{}

		// check path for index
		_, stringPath, err := makePath(path)
		if err != nil {
			return nil, err
		}

		// get working directory based on path
		objectDir,err := jmespath.Search(*stringPath, object)
		if err != nil {
			return nil, err
		}

		// marshal out to json
		rawObj, err := json.Marshal(objectDir)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(rawObj, kv)
		fmt.Println(child)
		fmt.Println(string(rawObj))

	}

	return &object, nil
}


// construct the string path from the path slice
// strips off index for last object to properly handle slices
func makePath(path []string) (*int, *string, error) {


	location := &path[len(path)-1]

	// compile regex to check for index in path
	index := regexp.MustCompile("^.*\\[[\\d]+\\]$")

	if index.MatchString(*location) {
		// regex to match index int and convert down to int from int64
		locationIndex := regexp.MustCompile("[\\d]+").FindString(*location)
		i64, err := strconv.ParseInt(locationIndex, 10, 8)
		if err != nil {
			return nil, nil, err
		}
		locationInt := int(i64)
		// regex to find string in path
		locationName := regexp.MustCompile("\\[[\\d]+\\]").ReplaceAllString(*location, "")

		path[len(path)-1] = locationName

		// combine the sliced path into jmespath format
		stringPath := parsing.PathFormatter(path)

		return &locationInt, &stringPath, nil
	}

	// combine the sliced path into jmespath format
	stringPath := parsing.PathFormatter(path)

	return nil, &stringPath, nil
}
















