package operator

import (
	"fmt"
	"github.com/beard1ess/gauss/parsing"
	"reflect"
	"strconv"
)

func keys(original parsing.Keyvalue, modified parsing.Keyvalue, path []string, objectDiff parsing.ConsumableDifference) parsing.ConsumableDifference {
	for k, v := range modified {
		if parsing.IndexOf(parsing.ListStripper(original), k) == -1 {
			added := parsing.AddedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
			objectDiff.Added = append(objectDiff.Added, added)
			delete(modified, k)
			fmt.Println("DELETED:  ", k)
		}
	}
	for k, v := range original {
		if parsing.IndexOf(parsing.ListStripper(modified), k) == -1 {
			removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
			objectDiff.Removed = append(objectDiff.Removed, removed)
			delete(original, k)
			fmt.Println("DELETED:  ", k)
		}
	}

	return recursion(original, modified, path, objectDiff)
}

func recursion(original parsing.Keyvalue, modified parsing.Keyvalue, input_path []string, objectDiff parsing.ConsumableDifference) parsing.ConsumableDifference {

	path := make([]string, len(input_path))
	copy(path, input_path)
	if reflect.DeepEqual(original, modified) {
		return objectDiff
	}

	if !(parsing.UnorderedKeyMatch(original, modified)) {

		objectDiff = keys(original, modified, path, objectDiff)

	}
	for k := range original {
		var valOrig, valMod interface{}
		if reflect.TypeOf(original).Kind() == reflect.String {
			valOrig = original
		} else {
			valOrig = original[k]
		}
		if reflect.TypeOf(modified).Kind() == reflect.String {
			valMod = modified
		} else {
			valMod = modified[k]
		}

		if !(reflect.DeepEqual(valMod, valOrig)) {
			if reflect.TypeOf(valOrig).Kind() == reflect.Map {
				path = append(path, k)
				recursion(parsing.Remarshal(valOrig), parsing.Remarshal(valMod), path, objectDiff)
				return objectDiff
			} else if reflect.TypeOf(valOrig).Kind() == reflect.Slice {
				var match bool
				valOrig, _ := valOrig.([]interface{})
				valMod, _ := valMod.([]interface{})
				path = append(path, k)
				npath := make([]string, len(path))
				if len(valOrig) != len(valMod) {
					if len(valOrig) > len(valMod) {
						for i := range valOrig {
							for ii := range valMod {
								if reflect.DeepEqual(valOrig[i], valMod[ii]) {
									match = true
								} else if i == ii {
									iter := len(path) - 1
									path[iter] = path[iter] + "[" + strconv.Itoa(i) + "]"
									recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]), path, objectDiff)
								}
							}
							if !(match) {
								removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path),
									Key: k, Value: valOrig}
								objectDiff.Removed = append(objectDiff.Removed, removed)
							} else {
								match = false
							}
						}

					} else {
						for i := range valMod {
							for ii := range valOrig {
								if reflect.DeepEqual(valOrig[ii], valMod[i]) {
									match = true
								} else if i == ii {
									iter := len(path) - 1
									path[iter] = path[iter] + "[" + strconv.Itoa(i) + "]"
									recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]), path, objectDiff)
								}
							}
							if !(match) {
								added := parsing.AddedDifference{Path: parsing.PathFormatter(path),
									Key: k, Value: valMod}
								objectDiff.Added = append(objectDiff.Added, added)
							} else {
								match = false
							}
						}
					}
				} else {
					for i := range valOrig {
						copy(npath, path)
						if !(reflect.DeepEqual(valOrig[i], valMod[i])) {
							iter := len(npath) - 1
							npath[iter] = npath[iter] + "[" + strconv.Itoa(i) + "]"
							recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]), npath, objectDiff)
						}
					}
					return objectDiff
				}
			} else {
				changed := parsing.ChangedDifference{Path: parsing.PathFormatter(path), Key: k,
					OldValue: valOrig, NewValue: valMod}
				objectDiff.Changed = append(objectDiff.Changed, changed)
				return objectDiff
			}
		}
	}

	return objectDiff
}

func Recursion(original parsing.Keyvalue, modified parsing.Keyvalue, path []string) parsing.ConsumableDifference {

	var objectDiff = parsing.ConsumableDifference{}

	return recursion(original, modified, path, objectDiff)
}
