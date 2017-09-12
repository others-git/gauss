package operator

import (
	"fmt"
	"github.com/beard1ess/gauss/parsing"
	"os"
	"reflect"
	"strconv"
)

func recursion(
	original parsing.Keyvalue,
	modified parsing.Keyvalue,
	path []string,
	objectDiff parsing.ConsumableDifference,
) parsing.ConsumableDifference {

	kListModified := parsing.ListStripper(modified)
	kListOriginal := parsing.ListStripper(original)
	if len(kListModified) > 1 || len(kListOriginal) > 1 {
		proc := true

		// If a key which exists in original doesn't exist in modified, add it
		// to the list of removed keys.
		for k, v := range original {
			if parsing.IndexOf(kListModified, k) == -1 {
				removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				objectDiff.Removed = append(objectDiff.Removed, removed)
				proc = false
			}
		}

		// If a key which exists in modified doesn't exist in original, add it
		// to the list of added keys.
		for k, v := range modified {
			if parsing.IndexOf(kListOriginal, k) == -1 {
				added := parsing.AddedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				objectDiff.Added = append(objectDiff.Added, added)
				proc = false
			}
		}

		if proc {
			for k := range original {
				recursion(parsing.Keyvalue{k: original[k]}, parsing.Keyvalue{k: modified[k]}, path, objectDiff)
			}
		}
		return objectDiff
	}

	// for each key in original
	for k := range original {
		var npath []string
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
				npath = append(path, k)
				recursion(parsing.Remarshal(valOrig), parsing.Remarshal(valMod), npath, objectDiff)
				return objectDiff
			} else if reflect.TypeOf(valOrig).Kind() == reflect.Slice {
				valOrig, _ := valOrig.([]interface{})
				valMod, _ := valMod.([]interface{})
				if len(valOrig) != len(valMod) {
					// TODO array length differences, how to interpret?
					fmt.Println("Cannot handle array length differences yet, sorry not sorry; kind of sorry.")
					os.Exit(1)
				} else {
					for i := range valOrig {
						if !(reflect.DeepEqual(valOrig[i], valMod[i])) {
							iter := len(path) - 1
							path[iter] = path[iter] + "[" + strconv.Itoa(i) + "]"
							recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]), path, objectDiff)
							return objectDiff
						}
					}
				}
			} else {
				changed := parsing.ChangedDifference{Path: parsing.PathFormatter(path), Key: k,
					OldValue: valOrig, NewValue: valMod}
				objectDiff.Changed = append(objectDiff.Changed, changed)
				return objectDiff
			}
		}
		return objectDiff
	}
	return objectDiff
}

func Recursion(original parsing.Keyvalue, modified parsing.Keyvalue, path []string) parsing.ConsumableDifference {

	objectDiff := parsing.ConsumableDifference{}

	return recursion(original, modified, path, objectDiff)
}
