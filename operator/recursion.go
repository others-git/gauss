package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"reflect"
	"strconv"
)

var ObjectDiff = parsing.ConsumableDifference{}

func recursion(original parsing.Keyvalue, modified parsing.Keyvalue, path []string) {
	if len(parsing.ListStripper(modified)) > 1 || len(parsing.ListStripper(original)) > 1 {
		proc := true
		for k, v := range original {
			if parsing.IndexOf(parsing.ListStripper(modified), k) == -1 {
				removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				ObjectDiff.Removed = append(ObjectDiff.Removed, removed)
				proc = false
			}
		}
		for k, v := range modified {
			if parsing.IndexOf(parsing.ListStripper(original), k) == -1 {
				added := parsing.AddedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				ObjectDiff.Added = append(ObjectDiff.Added, added)
				proc = false
			}
		}
		if proc {
			for k := range original {
				recursion(parsing.Keyvalue{k: original[k]}, parsing.Keyvalue{k: modified[k]}, path)
			}
		}

	} else {

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
					recursion(parsing.Remarshal(valOrig), parsing.Remarshal(valMod), path)
					return
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
										recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]), path)
									}
								}
								if !(match) {
									removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path),
										Key:                                   k, Value: valOrig}
									ObjectDiff.Removed = append(ObjectDiff.Removed, removed)
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
										recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]), path)
									}
								}
								if !(match) {
									added := parsing.AddedDifference{Path: parsing.PathFormatter(path),
										Key:                               k, Value: valMod}
									ObjectDiff.Added = append(ObjectDiff.Added, added)
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
								recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]), npath)
							}
						}
						return
					}
				} else {
					changed := parsing.ChangedDifference{Path: parsing.PathFormatter(path), Key: k,
						OldValue:                              valOrig, NewValue: valMod}
					ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
					return
				}
			}
		}
	}
}

func Recursion(original parsing.Keyvalue, modified parsing.Keyvalue, path []string) parsing.ConsumableDifference {
	recursion(original, modified, path)
	return ObjectDiff
}
