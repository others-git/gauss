package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"fmt"
	"github.com/jmespath/go-jmespath"
	"regexp"
	"strconv"
	"reflect"
	"encoding/json"
)

// https://github.com/golang/go/wiki/SliceTricks


// Patch: Creates a new object given a 'patch' and 'original'
func Patch(patch *parsing.ConsumableDifference, original *parsing.Gaussian) (*interface{}, error) {
	originalObject := &original.Data

	var newObject interface{}
	// Updated order Index > Changed > Removed > Added

	newObject, err := iterateAdded(patch.Added, *originalObject)
	if err != nil {
 	  return nil, err
   }

	newObject, err = iterateChanged(patch.Changed, *originalObject)
	if err != nil {
		return nil, err
	}


	res, _ := json.Marshal(newObject)
	fmt.Println(string(res))
	return &newObject, nil
}


// iterate over changed objects
func iterateAdded(added []parsing.AddedDifference, originalObject interface{}) (*interface{}, error) {
	var newObject interface{}

	// Iterate over added objects
	for _, i := range added {

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

		// wrap child object to create new object
		newObject, err = addParent(slicedPath, *childObject, originalObject)
		if err != nil {
			fmt.Println(err)
		}
	}

	return &newObject, nil
}



// iterate over changed objects
func iterateChanged(changed []parsing.ChangedDifference, originalObject interface{}) (*interface{}, error) {
	var newObject interface{}

	// Iterate over added objects
	for _, i := range changed {

		originPath := i.Path
		key := i.Key
		value := i.NewValue

		// validate jmespath
		_, err :=  jmespath.Compile(originPath)
		if err != nil {
			nErr := fmt.Errorf("failed to compile provided path: %T", err)
			return nil, nErr
		}

		// slice up path
		slicedPath := parsing.PathSplit(originPath)

		// create child object
		childObject, err := replaceChild(slicedPath, key, value, originalObject)
		if err != nil {
			return nil, err
		}

		// wrap child object to create new object
		newObject, err = addParent(slicedPath, *childObject, originalObject)
		if err != nil {
			fmt.Println(err)
		}
	}

	return &newObject, nil
}


// same as create but replaces slice index rather than inserting
func replaceChild(path []string, key string, value interface{}, object interface{}) (*interface{}, error) {

	var newObject interface{}

	// check path for index
	index, stringPath, err := makePath(path)
	if err != nil {
		return nil, err
	}

	// get working directory based on path
	objectDir,err := jmespath.Search(*stringPath, object)
	if err != nil {
		return nil, err
	}
	// determine what type of object we need to make - NEED MORE CHECKS
	if key != "" {
		// create k[v] type and return
		newChild := map[string]interface{}{
			key: value,
		}

		// reduce maps
		if reflect.TypeOf(objectDir).Kind() == reflect.Map {
			newObject = mapReduce(objectDir.(map[string]interface{}), newChild)
		} else {
			newObject = newChild
		}

	} else {

		newObject = value
	}

	// replace logic for slice value
	if reflect.TypeOf(objectDir).Kind() == reflect.Slice {
		//TODO: do thing with index
		// cast to slice of interfaces
		objectSlice := objectDir.([]interface{})

		// insert into slice
		objectSlice[*index] = newObject
		newObject = objectSlice
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
		return nil, err
	}
	// determine what type of object we need to make - NEED MORE CHECKS
	if key != "" {
		// create k[v] type and return
		newChild := map[string]interface{}{
			key: value,
		}
		// reduce maps
		if reflect.TypeOf(objectDir).Kind() == reflect.Map {
			newObject = mapReduce(objectDir.(map[string]interface{}), newChild)
		} else {
			newObject = newChild
		}

	} else {

		newObject = value
	}

	// Update logic for slice value - NEED MORE CHECKS
	if reflect.TypeOf(objectDir).Kind() == reflect.Slice {
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
	var objectDir interface{}

	// adjust path
	objectName := path[len(path)-1]
	lessPath := path[:len(path)-1]
	pathLen := len(path)-1

	if pathLen > 0 {
		// get string path
		_, stringPath, err := makePath(lessPath)
		if err != nil {
			return nil, err
		}

		// get working directory based on path
		objectDir, err = jmespath.Search(*stringPath, object)
		if err != nil {
			nErr := fmt.Errorf("invalid jmespath expression: %q", *stringPath)
			fmt.Println(objectName)
			fmt.Println(lessPath)
			fmt.Println(pathLen)
			return nil, nErr
		}
	} else {
		objectDir = object
	}

	//tmpObject := make(map[string]interface{})
	objectDir.(map[string]interface{})[objectName] = child

	if pathLen > 0 {
		return addParent(lessPath, objectDir, object)
	}
	return &objectDir, nil
}


// construct the string path from the path slice
// strips off index for last object to properly handle slices
func makePath(path []string) (*int, *string, error) {

	if len(path) == 0 {
		null := ""
		return nil, &null, nil
	}
	location := &path[len(path)-1]
	// remove any escaped quotes \"
	*location = regexp.MustCompile("\\\"").ReplaceAllString(*location, "")
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


func mapReduce(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}
	return a
}














