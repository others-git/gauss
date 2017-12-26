package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"reflect"
	"strconv"
	"regexp"
	"encoding/json"
)

type diff struct{
	skip *regexp.Regexp
}


func (d *diff) recursion(

	original interface{},
	modified interface{},
	path *[]string,
	ObjectDiff *parsing.ConsumableDifference,

) error {
	// check if we should skip value
	if d.skip != nil {
		res,err := json.Marshal(original)
		if err != nil {
			return err
		}
		res2,err := json.Marshal(modified)
		if err != nil {
			return err
		}
		if d.skip.MatchString(string(res)) || d.skip.MatchString(string(res2)) {
			return nil
		}
	}

	// equal values
	if reflect.DeepEqual(original, modified) {
		return nil
	}

	// grab object types
	originalType := reflect.TypeOf(original).Kind()
	modifiedType := reflect.TypeOf(modified).Kind()

	// type mismatch is a difference
	if originalType != modifiedType {
		changed := parsing.ChangedDifference{Path: parsing.CreatePath(*path),
			OldValue: original, NewValue: modified}
		ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
		return nil
	}

	// handle both values being a map
	if originalType == reflect.Map && modifiedType == reflect.Map {
		original := original.(map[string]interface{})
		modified := modified.(map[string]interface{})

		if !(parsing.UnorderedKeyMatch(original, modified)) {
			// check for key differences at the object's top level
			for k, v := range modified {

				if parsing.IndexOf(parsing.GetSliceOfKeys(original), k) == -1 {
					added := parsing.AddedDifference{Path: parsing.CreatePath(*path), Key: k, Value: v}
					ObjectDiff.Added = append(ObjectDiff.Added, added)
					delete(modified, k)
				}
			}
			for k, v := range original {

				if parsing.IndexOf(parsing.GetSliceOfKeys(modified), k) == -1 {
					removed := parsing.RemovedDifference{Path: parsing.CreatePath(*path), Key: k, Value: v}
					ObjectDiff.Removed = append(ObjectDiff.Removed, removed)
					delete(original, k)
				}
			}

			err := d.recursion(original, modified, path, ObjectDiff)
			if err != nil {
				return err
			}
			return nil

		} else if len(parsing.GetSliceOfKeys(original)) > 1 || len(parsing.GetSliceOfKeys(modified)) > 1 {
			// if there is more than 1 key, iterate through each and return
			for k := range original {
				err := d.recursion(map[string]interface{}{k: original[k]}, map[string]interface{}{k: modified[k]}, path, ObjectDiff)
				if err != nil {
					return err
				}
			}
			return nil
		}

	}
	// what gets into this area, strings and slices

	switch originalType{

	case reflect.Map:


		// cast to map
		original := original.(map[string]interface{})
		modified := modified.(map[string]interface{})

		err := d.mapHandler(original, modified, *path, ObjectDiff)
		if err != nil {
			return err
		}

	case reflect.Slice:

		// cast to slice
		original := original.([]interface{})
		modified := modified.([]interface{})

		// pass slices off to handler
		err := d.sliceHandler(original, modified, *path, ObjectDiff)
		if err != nil {
			return err
		}
	case reflect.String:
		// if type is string, escape non printable characters
		original := original.(string)
		modified := modified.(string)
		original = strconv.Quote(original)
		modified = strconv.Quote(modified)

		changed := parsing.ChangedDifference{Path: parsing.CreatePath(*path),
			OldValue: original, NewValue: modified}
		ObjectDiff.Changed = append(ObjectDiff.Changed, changed)

	default:

		changed := parsing.ChangedDifference{Path: parsing.CreatePath(*path),
			OldValue: original, NewValue: modified}
		ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
		/*

		err := fmt.Errorf("unknown type error, please report as bug." +
			"\noriginal type: %s \nmodified type: %s\n=====================" +
				"\noriginal value: %s \nmodified value: %s",
				originalType, modifiedType, original, modified)


		return err
		*/
	}

	return nil

}


func (d *diff) mapHandler(

	original map[string]interface{},
	modified map[string]interface{},
	path []string,
	diff *parsing.ConsumableDifference,

) error {

	for k := range original {

		originalValue := original[k]
		modifiedValue := modified[k]

		// type mismatch
		if reflect.TypeOf(originalValue) != reflect.TypeOf(modifiedValue) {
			changed := parsing.ChangedDifference{Path: parsing.CreatePath(path),
				Key: k, OldValue: originalValue, NewValue: modifiedValue}
			diff.Changed = append(diff.Changed, changed)
			return nil
		// maps
		} else if reflect.TypeOf(originalValue).Kind() == reflect.String {
			changed := parsing.ChangedDifference{Path: parsing.CreatePath(path),
				Key: k, OldValue: originalValue, NewValue: modifiedValue}
			diff.Changed = append(diff.Changed, changed)
			return nil
		} else {
			// Update the working path
			path = append(path, k)
			err := d.recursion(originalValue, modifiedValue, &path, diff)
			if err != nil {
				return err
			}
			return nil
			// Slice handler
		}
	}

	return nil
}



func (d *diff) sliceHandler(

	original []interface{},
	modified []interface{},
	path []string,
	diff *parsing.ConsumableDifference,

) error {

	originalLength := len(original)
	modifiedLength := len(modified)

	// handle length mismatch
	if originalLength != modifiedLength {

		if originalLength > modifiedLength {
			// handle multiple length differences
			lengthDifference := originalLength - modifiedLength
			for i := 1 ; i <= lengthDifference; i++ {
				index := originalLength-i
				// if original is longer we know an item was removed
				removed := parsing.RemovedDifference{Path: parsing.CreatePath(parsing.SliceIndex(index, path)),
					Value: original[index]}
				diff.Removed = append(diff.Removed, removed)

				// remove what we parsed
				original = append(original[:index], original[index+1:]...)
			}

		} else {
			// handle multiple length differences
			lengthDifference := modifiedLength - originalLength
			for i := 1 ; i <= lengthDifference; i++ {
				index := modifiedLength-i
				// if modified is longer we know an item was added
				added := parsing.AddedDifference{Path: parsing.CreatePath(parsing.SliceIndex(index, path)),
					Value: modified[index]}
				diff.Added = append(diff.Added, added)

				// remove what we parsed
				modified = append(modified[:index], modified[index+1:]...)
			}

		}
	}

	// if length are the same iterate over all and recurse
	for i := range original {
		path := parsing.SliceIndex(i, path)
		err := d.recursion(original[i], modified[i], &path, diff)
		if err != nil {
			return err
		}
	}

	return nil
}


// Recursion iterate over objects to find differences
func Recursion(

	original interface{},
	modified interface{},
	path []string,
	regSkip *regexp.Regexp,

) (*parsing.ConsumableDifference, error) {

	var differ diff

	var ObjectDiff parsing.ConsumableDifference
	err := differ.recursion(original, modified, &path, &ObjectDiff)
	if err != nil {
		return nil, err
	}
	return &ObjectDiff, nil
}

