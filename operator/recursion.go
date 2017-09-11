package operator

import (
	"strconv"
	"fmt"
	"os"
	"github.com/beard1ess/gauss/parsing"
	"reflect"
)

var ObjectDiff = parsing.ConsumableDifference{}

func recursion(original parsing.Keyvalue, modified parsing.Keyvalue, path []string) {
	kListModified := parsing.ListStripper(modified)
	kListOriginal := parsing.ListStripper(original)
	if len(kListModified) > 1 || len(kListOriginal) > 1 {
		proc := true
		for k, v := range original {
			if parsing.IndexOf(kListModified, k) == -1 {
				removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				ObjectDiff.Removed = append(ObjectDiff.Removed, removed)
				proc = false
			}
		}
		for k, v := range modified {
			if parsing.IndexOf(kListOriginal, k) == -1 {
				added := parsing.AddedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				ObjectDiff.Added = append(ObjectDiff.Added, added)
				proc = false
			}
		}
		if proc {
			for k := range original {
				Recursion(parsing.Keyvalue{k:original[k]},parsing.Keyvalue{k:modified[k]},path)
			}
		}
		return
	}
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
				Recursion(parsing.Remarshal(valOrig), parsing.Remarshal(valMod), npath)
				return
			} else if reflect.TypeOf(valOrig).Kind() == reflect.Slice {
				valOrig,_ := valOrig.([]interface{})
				valMod,_ := valMod.([]interface{})
				if len(valOrig) != len(valMod) {
					// TODO array length differences
					fmt.Println("Cannot handle array length differences yet, sorry not sorry; kind of sorry.")
					os.Exit(1)
				} else {
					for i := range valOrig {
						if !(reflect.DeepEqual(valMod[i], valOrig[i])) {
							path[len(path)-1] = path[len(path)-1] + "[" + strconv.Itoa(i) + "]"
							changed := parsing.ChangedDifference{Path: parsing.PathFormatter(path), Key: k,
								OldValue: valOrig[i], NewValue: valMod[i]}
							ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
							return
						}
					}
				}
			} else {
				changed := parsing.ChangedDifference{Path: parsing.PathFormatter(path), Key: k,
					OldValue: valOrig, NewValue: valMod}
				ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
				return
			}
		}
		return
	}
	return
}

func Recursion(original parsing.Keyvalue, modified parsing.Keyvalue, path []string) parsing.ConsumableDifference {
	recursion(original, modified, path)
	return ObjectDiff
}