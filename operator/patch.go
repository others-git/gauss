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


type Patched struct{
	NewObject interface{}
	Patch *parsing.ConsumableDifference
	Skip *regexp.Regexp
}

// Patch Creates a new object given a 'patch' and 'original'
func Patch(patch *parsing.ConsumableDifference, original *parsing.Gaussian, regSkip *regexp.Regexp) (*interface{}, error) {
	var newObject interface{}
	P := Patched{newObject, patch, regSkip}
	P.NewObject = original.Data
	P.patch()

	return &P.NewObject, nil
}

// Patch Creates a new object given a 'patch' and 'original'
func (p *Patched) patch() (*interface{}, error) {

	var err error

	// remove
	err = p.iterateRemoved()
	if err != nil {
		return nil, err
	}

	// add
	err = p.iterateAdded()
	if err != nil {
		return nil, err
	}
	
	// change
	err = p.iterateChanged()
	if err != nil {
		return nil, err
	}

	return &p.NewObject, nil
}

// iterate over objects to remove
func (p *Patched) iterateRemoved() error {
	Removing:
	// iterate over removed objects
	for _, i := range p.Patch.Removed {

		originPath := i.Path
		key := i.Key
		value := i.Value

		// skip matched regex
		if p.Skip != nil {
			res,err := json.Marshal(value)
			if err != nil {
				return err
			}
			if p.Skip.MatchString(string(res)) {
				continue Removing
			}
		}

		// slice up path
		slicedPath := parsing.PathSplit(originPath)

		// create child object
		childObject, err := p.removeChild(slicedPath, key, value, p.NewObject)
		if err != nil {
			return err
		}

		_, trimPath := trimIndex(slicedPath)
		// wrap child object to create new object
		_, err = p.addParent(trimPath, *childObject, p.NewObject)
		if err != nil {
			return err
		}
	}

	return nil
}

// iterate over objects to add
func (p *Patched) iterateAdded() error {
	Adding:
	// Iterate over added objects
	for _, i := range p.Patch.Added {

		originPath := i.Path
		key := i.Key
		value := i.Value

		// skip matched regex
		if p.Skip != nil {

			res,err := json.Marshal(value)
			if err != nil {
				return err
			}

			if p.Skip.MatchString(string(res)) {
				continue Adding
			}
		}

		// slice up path
		slicedPath := parsing.PathSplit(originPath)

		// create child object
		childObject, err := p.createChild(slicedPath, key, value, p.NewObject)
		if err != nil {
			return err
		}

		_, trimPath := trimIndex(slicedPath)
		// wrap child object to create new object
		_, err = p.addParent(trimPath, *childObject, p.NewObject)
		if err != nil {
			return err
		}
	}

	return nil
}

// iterate over objects to change
func (p *Patched) iterateChanged() error {
	Changing:
	// Iterate over added objects
	for _, i := range p.Patch.Changed {

		originPath := i.Path
		key := i.Key
		value := i.NewValue

		// skip matched regex
		if p.Skip != nil {
			res,err := json.Marshal(value)
			if err != nil {
				return err
			}
			res2,err := json.Marshal(i.OldValue)
			if err != nil {
				return err
			}
			if p.Skip.MatchString(string(res)) || p.Skip.MatchString(string(res2)) {
				continue Changing
			}
		}

		// skip regexp

		// slice up path
		slicedPath := parsing.PathSplit(originPath)


		// create child object
		childObject, err := p.replaceChild(slicedPath, key, value, p.NewObject)
		if err != nil {
			return err
		}

		_, trimPath := trimIndex(slicedPath)
		// wrap child object to create new object
		_, err = p.addParent(trimPath, *childObject, p.NewObject)
		if err != nil {
			return err
		}

	}

	return nil
}

////////

func (p *Patched) removeChild(path []string, key string, value interface{}, object interface{}) (*interface{}, error) {

	var orphan interface{}
	var newObject interface{}

	objectName := path[len(path)-1]

	index,tmp := trimIndex(path)

	// get working directory based on path
	usePath := parsing.CreatePath(tmp)

	//objectDir, err := compiledPath.Search(object)
	objectDir, err := jmespath.Search(usePath, object)
	if err != nil || objectDir == nil {
		nErr := fmt.Errorf("\n::::::::::::::::::::::::::::::::::::::\n" +
			"\nerror: %v\nraw path: %v\n" +
			"result: %v\nstack: %s\n\n" +
			"::::::::::::::::::::::::::::::::::::::\n", err, strconv.Quote(usePath), objectDir, object)
		return nil, nErr
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
			objectSlice = append(objectSlice[:*index], objectSlice[*index+1:]...)
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
func (p *Patched) replaceChild(path []string, key string, value interface{}, object interface{}) (*interface{}, error) {

	var newObject interface{}
	var err error

	valueType := reflect.TypeOf(value).String()
	switch valueType{
	case "string":
		val := value.(string)

		value,err = strconv.Unquote(val)
		if err != nil {
			value = val
		}
	default:
	}
	/*

	_,tmp := trimIndex(path)

	// create index and jmespath
	index, _, compiledPath, err := makePath(tmp)
	if err != nil {
		return nil, err
	}
	*/

	// check path for index
	index,_, compiledPath, err := makePath(path)
	if err != nil {
		return nil, err
	}

	objectDir, err := compiledPath.Search(object)
	tmpObject := objectDir
	if err != nil || objectDir == nil {
		nErr := fmt.Errorf("\n::::::::::::::::::::::::::::::::::::::\n" +
				"\nerror: %v\nraw path: %v\n" +
				"compiled path: %v\nresult: %v\nstack: %s\n\n" +
				"::::::::::::::::::::::::::::::::::::::\n", err, path, *compiledPath, objectDir, object)
		return nil, nErr
	}

	// determine what type of object we need to make - NEED MORE CHECKS
	if key != "" {
		// step into slice if valid
		if index != nil {
			objectDir = objectDir.([]interface{})[*index]
		}

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
	if index != nil {
		tmpObject.([]interface{})[*index] = newObject
		newObject = tmpObject
	}

	return &newObject, nil
}

// creates new child object from key and value
func (p *Patched) createChild(path []string, key string, value interface{}, object interface{}) (*interface{}, error) {

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

	objectDir, err := compiledPath.Search(object)
	tmpObject := objectDir
	if err != nil || objectDir == nil {
		nErr := fmt.Errorf("\n::::::::::::::::::::::::::::::::::::::\n" +
			"\npath expression returned nil\nquery path: %v\nresult: %v\n\n" +
			"::::::::::::::::::::::::::::::::::::::\n", *compiledPath, objectDir)
		return nil, nErr
	}

	// determine what type of object we need to make - NEED MORE CHECKS
	if key != "" {

		// step into slice if valid
		if index != nil {
			objectDir = objectDir.([]interface{})[*index]
		}

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

	if index != nil {
		tmpObject.([]interface{})[*index] = newObject
		newObject = tmpObject
	}

	return &newObject, nil
}

// recreate the template with the updated child object
func (p *Patched) addParent(path []string, lastObject interface{}, stack interface{}) (*interface{}, error) {
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
	multiIndexReg := regexp.MustCompile("\\[[\\d]\\]\\[[\\d]\\]+")

	if 	multiIndexReg.MatchString(lastItem) {
		var err error

		// parse all nested slices
		child, err = p.nestedSlice(path, lastObject, stack)
		if err != nil {
			return nil, err
		}
		// replace indexes to set as value
		lastItem = regexp.MustCompile("\\[[\\d]\\]+").ReplaceAllString(lastItem, "")
		lessPath = append(lessPath, lastItem)

		return p.addParent(lessPath, child, stack)
	} else {
		child = lastObject
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
				"\nerror reconstructing lastObject body\npath expression returned nil\nquery path: %v\nresult: %v\n\n" +
					"::::::::::::::::::::::::::::::::::::::\n", *compiledPath, newObject)
			return nil, nErr
		}
	} else {
		newObject = stack
	}

	// handle slice within the path
	switch reflect.TypeOf(newObject).String() {

	case "slice":

		if index == nil {
			nErr := fmt.Errorf("operating stack is type slice but path does not contain index")
			return nil, nErr
		}
		newObject.([]interface{})[*index] = child

	default:

		if regexp.MustCompile("\\[[\\d]\\]").MatchString(lastItem) {

			// if last item in path has an index value we have to unwrap and update the slice underneath
			tmp := []string{lastItem}
			indx, tmp := trimIndex(tmp)
			tmpObj := newObject.(map[string]interface{})[tmp[0]]


			tmpObj.([]interface{})[*indx] = child
			newObject.(map[string]interface{})[tmp[0]] = tmpObj

		} else {

			// remove any quoted values from the key name
			objectName,err := strconv.Unquote(lastItem)
			if err != nil {
				newObject.(map[string]interface{})[lastItem] = child
			} else {
				newObject.(map[string]interface{})[objectName] = child
			}
		}
	}

	if pathLen > 0 {
		return p.addParent(lessPath, newObject, stack)
	}

	p.NewObject = newObject

	return &newObject, nil
}

// if a path has nested slices, eg key[1][2] recurse over them
func (p *Patched) nestedSlice(path []string, child interface{}, stack interface{}) (*interface{}, error) {


	_, tmpPath := trimIndex(path)
	index, _, compiledPath, _:= makePath(tmpPath)
	lastItem := tmpPath[len(tmpPath)-1]


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
		return p.nestedSlice(tmpPath, objectSlice, stack)
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


	location = regexp.MustCompile("\\[[\\d]\\]+$").ReplaceAllString(location, "")

	index, tmp := trimIndex(clone)

	// combine the sliced path into jmespath format
	stringPath := parsing.CreatePath(tmp)

	// validate jmespath
	compiled, err :=  jmespath.Compile(stringPath)
	if err != nil {
		nErr := fmt.Errorf("failed to compile provided path: %T\npath len: %v", err, len(stringPath))
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
func trimIndex(path []string) (*int, []string) {

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

		return &index, tmp
	}


	return nil, tmp
}









