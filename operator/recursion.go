package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"reflect"
	"strconv"
	"fmt"
)

func recursion(

	original parsing.Keyvalue,
	modified parsing.Keyvalue,
	input_path []string,
	ObjectDiff parsing.ConsumableDifference,

) parsing.ConsumableDifference {

	path := make([]string, len(input_path))
	copy(path, input_path)
	if reflect.DeepEqual(original, modified) {
		return ObjectDiff
	}

	if !(parsing.UnorderedKeyMatch(original, modified)) {

		for k, v := range modified {
			if parsing.IndexOf(parsing.ListStripper(original), k) == -1 {
				added := parsing.AddedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				ObjectDiff.Added = append(ObjectDiff.Added, added)
				delete(modified, k)
			}
		}
		for k, v := range original {
			if parsing.IndexOf(parsing.ListStripper(modified), k) == -1 {
				removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				ObjectDiff.Removed = append(ObjectDiff.Removed, removed)
				delete(original, k)
			}
		}

		ObjectDiff = recursion(original, modified, path, ObjectDiff)
		return ObjectDiff

	} else if len(parsing.ListStripper(original)) > 1 || len(parsing.ListStripper(modified)) > 1 {

		for k := range original {
			ObjectDiff = recursion(
				parsing.Keyvalue{k: original[k]},
				parsing.Keyvalue{k: modified[k]},
				path, ObjectDiff)
		}
		return ObjectDiff
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
				// Specifically handle type mismatch
				if reflect.TypeOf(valOrig) != reflect.TypeOf(valMod) {
					changed := parsing.ChangedDifference{Path: parsing.PathFormatter(path),
						Key: k, OldValue: valOrig, NewValue: valMod}
					ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
					return ObjectDiff
				// Map handler
				} else if reflect.TypeOf(valOrig).Kind() == reflect.Map {
					// Update the working path
					path = append(path, k)
					ObjectDiff = recursion(parsing.Remarshal(valOrig), parsing.Remarshal(valMod), path, ObjectDiff)
					return ObjectDiff
				// Slice handler
				} else if reflect.TypeOf(valOrig).Kind() == reflect.Slice {

					// Variable setup
					var match bool
					valOrig, _ := valOrig.([]interface{})
					valMod, _ := valMod.([]interface{})
					// Update the working path and copy into a new var
					path = append(path, k)
					npath := make([]string, len(path))
					copy(npath, path)
					if len(valOrig) != len(valMod) {
						// If slice length mismatches we need to handle that a particular way
						if len(valOrig) > len(valMod) {
							for i := range valOrig {
								for ii := range valMod {
									if reflect.DeepEqual(valOrig[i], valMod[ii]) {

										match = true

									} else if i == ii {

										iter := len(path) - 1
										path[iter] = path[iter] + "[" + strconv.Itoa(i) + "]"
										ObjectDiff = recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]),
											path, ObjectDiff)

									}
								}
								if !(match) {
									removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path),
										Key: k, Value: valOrig}
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
										ObjectDiff = recursion(parsing.Remarshal(valOrig[i]),
											parsing.Remarshal(valMod[i]), path, ObjectDiff)
									}
								}
								if !(match) {
									added := parsing.AddedDifference{Path: parsing.PathFormatter(path),
										Key: k, Value: valMod}
									ObjectDiff.Added = append(ObjectDiff.Added, added)

								} else {
									match = false
								}
							}

						}
					} else {
						// If both slice lengths are equal
						for i := range valOrig {
							if !(reflect.DeepEqual(valOrig[i], valMod[i])) {
								iter := len(npath) - 1
								npath[iter] = npath[iter] + "[" + strconv.Itoa(i) + "]"
								ObjectDiff = recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]),
									npath, ObjectDiff)
							}
						}

					}
				} else {
					changed := parsing.ChangedDifference{Path: parsing.PathFormatter(path),
						Key: k, OldValue: valOrig, NewValue: valMod}
					ObjectDiff.Changed = append(ObjectDiff.Changed, changed)

				}
			}
		}
		return ObjectDiff
	}
}

func Recursion(original parsing.Keyvalue, modified parsing.Keyvalue, path []string) parsing.ConsumableDifference {
	var ObjectDiff = parsing.ConsumableDifference{}
	return recursion(original, modified, path, ObjectDiff)
}
