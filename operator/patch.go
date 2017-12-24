package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"fmt"
	"github.com/jmespath/go-jmespath"
	"regexp"
	"strconv"
	"reflect"
)

// https://github.com/golang/go/wiki/SliceTricks


// Patch Creates a new object given a 'patch' and 'original'
func Patch(patch *parsing.ConsumableDifference, original *parsing.Gaussian) (*interface{}, error) {
	originalObject := &original.Data

	var newObject interface{}
	// Updated order Removed > Added > Changed > Indexes

	// remove
	newObject, err := iterateRemoved(patch.Removed, *originalObject)
	if err != nil {
		return nil, err
	}

	// add
	newObject, err = iterateAdded(patch.Added, *originalObject)
	if err != nil {
		return nil, err
	}

	// alter
	newObject, err = iterateChanged(patch.Changed, *originalObject)
	if err != nil {
		return nil, err
	}

	return &newObject, nil
}

// iterate over objects to remove
func iterateRemoved(removed []parsing.RemovedDifference, originalObject interface{}) (*interface{}, error) {
	var newObject interface{}

	// iterate over removed objects
	for _, i := range removed {

		originPath := i.Path
		key := i.Key
		value := i.Value

		// validate jmespath
		_, err :=  jmespath.Compile(originPath)
		if err != nil {
			nErr := fmt.Errorf("failed to compile provided path: %T\npath: %s", err, originPath)
			return nil, nErr
		}

		// slice up path
		slicedPath := parsing.PathSplit(originPath)

		// create child object
		childObject, err := removeChild(slicedPath, key, value, originalObject)
		if err != nil {
			return nil, err
		}

		// wrap child object to create new object
		newObject, err = addParent(slicedPath, *childObject, originalObject)
		if err != nil {
			return nil, err
		}

	}

	return &newObject, nil
}

// iterate over objects to add
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
			nErr := fmt.Errorf("failed to compile provided path: %T\npath: %s", err, originPath)
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
			return nil, err
		}
	}

	return &newObject, nil
}

// iterate over objects to change
func iterateChanged(changed []parsing.ChangedDifference, originalObject interface{}) (*interface{}, error) {
	var newObject interface{}

	// Iterate over added objects
	for _, i := range changed {

		originPath := i.Path
		key := i.Key
		value := i.NewValue

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
			return nil, err
		}

	}

	return &newObject, nil
}

////////

func removeChild(path []string, key string, value interface{}, object interface{}) (*interface{}, error) {

	var orphan interface{}
	var newObject interface{}

	// check path for index
	index, _, compiledPath, err := makePath(path)
	if err != nil {
		return nil, err
	}

	objectName := path[len(path)-1]

	// get working directory based on path
	objectDir, err := compiledPath.Search(object)
	if err != nil {
		return nil, err
	}

	// determine what type of object to remove
	if key != "" {
		// create k[v] type and return
		orphan = map[string]interface{}{
			key: value,
		}
	} else {
		orphan = value
	}

	// replace logic for slice value
	if reflect.TypeOf(objectDir).Kind() == reflect.Slice {
		// cast to slice of interfaces
		objectSlice := objectDir.([]interface{})
		// if the object to orphan equals what's in the original object, drop it
		if reflect.DeepEqual(objectSlice[*index], orphan) {
			objectSlice = append(objectSlice[:*index], objectSlice[*index+1:])
		} else {
			nErr := fmt.Errorf("object to remove: %s\ndoes not match existing: %s",
				orphan, objectSlice[*index])
			return nil, nErr
		}
		newObject = objectSlice
	} else {

		// Cast to maps
		objectMap := objectDir.(map[string]interface{})
		orphan := orphan.(map[string]interface{})

		if parsing.MapMatchAny(orphan, objectMap) {
			delete(objectMap, key)
		} else {
			nErr := fmt.Errorf("object to remove: %s\ndoes not match existing: %s",
				orphan, objectMap[objectName])
			return nil, nErr
		}
		newObject = objectMap
	}

	return &newObject, nil
}

// same as create but replaces slice index rather than inserting
func replaceChild(path []string, key string, value interface{}, object interface{}) (*interface{}, error) {

	var newObject interface{}
	var err error

	valueType := reflect.TypeOf(value).Kind()
	switch valueType{
	case reflect.String:
		val := value.(string)

		value,err = strconv.Unquote(val)
		if err != nil {
			value = val
		}
	}


	// create index and jmespath
	index, _, compiledPath, err := makePath(path)
	if err != nil {
		return nil, err
	}

	// get working directory based on path

	objectDir, err := compiledPath.Search(object)
	if err != nil || objectDir == nil {
		nErr := fmt.Errorf("\n::::::::::::::::::::::::::::::::::::::\n" +
				"\nerror replacing value\n" +
				"path expression returned nil\nraw path: %v\n" +
				"compiled path: %v\nresult: %v\n\n" +
				"::::::::::::::::::::::::::::::::::::::\n", path, *compiledPath, objectDir)
		return nil, nErr
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
		// cast the object since it's a slice
		objectDir := objectDir.([]interface{})

		// create new slice
		objectSlice := make([]interface{}, len(objectDir))
		copy(objectSlice, objectDir)

		// insert into slice

		objectSlice[*index] = newObject

		newObject = objectSlice
	}

	return &newObject, nil
}

// creates new child object from key and value
func createChild(path []string, key string, value interface{}, object interface{}) (*interface{}, error) {

	var newObject interface{}
	var err error

	valueType := reflect.TypeOf(value).Kind()
	switch valueType{
	case reflect.String:
		val := value.(string)

		value,err = strconv.Unquote(val)
		if err != nil {
			value = val
		}
	default:

	}

	// check path for index
	index,_, compiledPath, err := makePath(path)
	if err != nil {
		return nil, err
	}

	// get working directory based on path
	objectDir, err := compiledPath.Search(object)
	if err != nil || objectDir == nil {
		nErr := fmt.Errorf("\n::::::::::::::::::::::::::::::::::::::\n" +
			"\npath expression returned nil\nquery path: %v\nresult: %v\n\n" +
			"::::::::::::::::::::::::::::::::::::::\n", *compiledPath, objectDir)
		return nil, nErr
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

		// if index is greater then total length we can't insert to add
		if *index > len(objectSlice) {
			// create new slice of index length and copy into it
			newSlice := make([]interface{}, *index)

			copy(newSlice, objectSlice)

			newSlice[*index] = newObject

			newObject = newSlice
		} else {
			// insert into slice
			objectSlice = append(objectSlice, 0)
			copy(objectSlice[*index+1:], objectSlice[*index:])
			objectSlice[*index] = newObject
			newObject = objectSlice
		}

	}

	return &newObject, nil
}

// recreate the template with the updated child object
func addParent(path []string, object interface{}, stack interface{}) (*interface{}, error) {
	var newObject interface{}
	var child interface{}


	// copy path into tmp to not mess with path
	tmp := make([]string, len(path))
	copy(tmp, path)

	// get various items
	lastItem := tmp[len(tmp)-1]
	lessPath := tmp[:len(tmp)-1]
	pathLen := len(tmp)-1

	// check if there is a multi index item in the path
	multIndexReg := regexp.MustCompile("\\[[\\d]\\]\\[[\\d]\\]+")
	allIndexReg := regexp.MustCompile("\\[[\\d]\\]+")
	if 	multIndexReg.MatchString(lastItem) {
		var err error

		// parse all nested slices
		child, err = nestedSlice(path, object, stack)
		if err != nil {
			return nil, err
		}
		// replace indexes to set as value
		lastItem = allIndexReg.ReplaceAllString(lastItem, "")
		lessPath = append(lessPath, lastItem)

		return addParent(lessPath, child, stack)
	} else {
		child = object
	}

	// pull parse the sliced path for some goodies
	index, _, compiledPath, err := makePath(lessPath)
	if err != nil {
		return nil, err
	}

	if pathLen > 0 {
		// get working directory based on path
		newObject, err = compiledPath.Search(stack)
		if err != nil || newObject == nil {
			nErr := fmt.Errorf("\n::::::::::::::::::::::::::::::::::::::\n" +
				"\nerror reconstructing object body\npath expression returned nil\nquery path: %v\nresult: %v\n\n" +
					"::::::::::::::::::::::::::::::::::::::\n", *compiledPath, newObject)
			return nil, nErr
		}
	} else {
		newObject = stack
	}

	// handle slice within the path
	if reflect.TypeOf(newObject).Kind() == reflect.Slice {
		fmt.Println("is slice")
		if index == nil {
			nErr := fmt.Errorf("operating stack is type slice but path does not contain index")
			return nil, nErr
		}
		newObject.([]interface{})[*index] = child
	}	else {

		// remove any quoted values from the key name
		objectName,err := strconv.Unquote(lastItem)
		if err != nil {
			//fmt.Printf("\n=====\nkey name: %s\noriginal name?: %s" +
			//	"\nerror: %s\n=====\n\n", objectName, lastItem,err)
			newObject.(map[string]interface{})[lastItem] = child
		} else {
			newObject.(map[string]interface{})[objectName] = child
		}
	}

	if pathLen > 0 {
		return addParent(lessPath, newObject, stack)
	}
	return &newObject, nil
}

// if a path has nested slices, eg key[1][2] recurse over them
func nestedSlice(path []string, child interface{}, stack interface{}) (*interface{}, error) {


	_, tmpPath := trimIndex(path)
	index, _, compiledPath, _:= makePath(*tmpPath)
	cln := *tmpPath
	lastItem := cln[len(cln)-1]


	// get working directory based on path
	objectDir, err := compiledPath.Search(stack)
	if err != nil || objectDir == nil {
		nErr := fmt.Errorf("\n::::::::::::::::::::::::::::::::::::::\n" +
			"\nerror reconstructing object body\npath expression returned nil\nquery path: %v\nresult: %v\n\n" +
			"::::::::::::::::::::::::::::::::::::::\n", *compiledPath, objectDir)
		return nil, nErr
	}

	objectSlice := objectDir.([]interface{})
	objectSlice[*index] = child

	multIndexReg := regexp.MustCompile("^\"?.*\"?\\[[\\d]]\\[[\\d]\\]+")
	if 	multIndexReg.MatchString(lastItem) {
		return nestedSlice(*tmpPath, objectSlice, stack)
	}

	return &objectDir, nil
}

// construct the string path from the sliced path
// strips off index for last object to properly handle slices
func makePath(path []string) (*int, *string, *jmespath.JMESPath, error) {

	// return nil if we need if here with no path
	if len(path) == 0 {
		return nil, nil, nil, nil
	}

	// copy path into tmp to not mess with path
	clone := make([]string, len(path))
	copy(clone, path)

	location := clone[len(clone)-1]
	// remove any escaped quotes \" - TODO: check test
	//*location = regexp.MustCompile("\\\"").ReplaceAllString(*location, "")


	index, tmp := trimIndex(clone)

	// combine the sliced path into jmespath format
	stringPath := parsing.CreatePath(*tmp)

	// validate jmespath
	compiled, err :=  jmespath.Compile(stringPath)
	if err != nil {
		nErr := fmt.Errorf("failed to compile provided path: %T\npath: %s", err, stringPath)
		return nil, nil, nil, nErr
	}


	return index, &location, compiled, nil
}

// reduce two maps
func mapReduce(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}
	return a
}

// trim last index off item and return the int
func trimIndex(path []string) (*int, *[]string) {

	// copy path into tmp to not mess with path
	tmp := make([]string, len(path))
	copy(tmp, path)

	location := tmp[len(tmp)-1]

	// compile regex to check for index on last item
	indexReg := regexp.MustCompile("^.*\\[[\\d]\\]$")
	if indexReg.MatchString(location) {

		// regex to match index int and convert down to int from int64
		locationIndex := regexp.MustCompile("[\\d]").FindString(location)
		i64, _ := strconv.ParseInt(locationIndex, 10, 8)

		index := int(i64)

		// regex to find string in path
		locationName := regexp.MustCompile("\\[[\\d]\\]$").ReplaceAllString(location, "")

		tmp[len(tmp)-1] = locationName

		return &index, &tmp
	}


	return nil, &tmp
}









