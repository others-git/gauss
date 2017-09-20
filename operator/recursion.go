package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"reflect"
)

func recursion(

	original parsing.Keyvalue,
	modified parsing.Keyvalue,
	path []string,
	ObjectDiff parsing.ConsumableDifference,

) parsing.ConsumableDifference {

	if reflect.DeepEqual(original, modified) {
		return ObjectDiff
	}

	if !(parsing.UnorderedKeyMatch(original, modified)) {

		for k, v := range modified {
			if parsing.IndexOf(parsing.Slicer(original), k) == -1 {
				added := parsing.AddedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				ObjectDiff.Added = append(ObjectDiff.Added, added)
				delete(modified, k)
			}
		}
		for k, v := range original {
			if parsing.IndexOf(parsing.Slicer(modified), k) == -1 {
				removed := parsing.RemovedDifference{Path: parsing.PathFormatter(path), Key: k, Value: v}
				ObjectDiff.Removed = append(ObjectDiff.Removed, removed)
				delete(original, k)
			}
		}

		ObjectDiff = recursion(original, modified, path, ObjectDiff)
		return ObjectDiff

	} else if len(parsing.Slicer(original)) > 1 || len(parsing.Slicer(modified)) > 1 {

		for k := range original {
			ObjectDiff = recursion(parsing.Keyvalue{k: original[k]}, parsing.Keyvalue{k: modified[k]}, path, ObjectDiff)
		}
		return ObjectDiff
	} else {

		for k := range original {
			valOrig := original[k]
			valMod := modified[k]

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
					valOrig, _ := valOrig.([]interface{})
					valMod, _ := valMod.([]interface{})
					// Update the working path and copy into a new var
					path = append(path, k)
					if len(valOrig) != len(valMod) {
						// If slice length mismatches we need to handle that a particular way
						if len(valOrig) > len(valMod) {
							for i := range valOrig {
							Mod:
								for ii := range valMod {
									if i == ii && reflect.DeepEqual(valOrig[i], valMod[ii]) {
										break Mod
									} else {
										if i != ii && reflect.DeepEqual(valOrig[i], valMod[ii]) {
											indexed := parsing.IndexDifference{OldIndex: i, NewIndex: ii, Value: valOrig[ii],
												Path: parsing.PathFormatter(path)}
											ObjectDiff.Indexes = append(ObjectDiff.Indexes, indexed)
											break Mod
										} else if i == ii && !(reflect.DeepEqual(valOrig[i], valMod[ii])) {
											if reflect.TypeOf(valOrig[i]).Kind() == reflect.String || reflect.TypeOf(valMod[ii]).Kind() == reflect.String {
												changed := parsing.ChangedDifference{Path: parsing.PathFormatter(parsing.PathSlice(i, path)),
													OldValue: valOrig[i], NewValue: valMod[i]}
												ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
												break Mod
											} else {
												ObjectDiff = recursion(parsing.Remarshal(valOrig[i]),
													parsing.Remarshal(valMod[i]), parsing.PathSlice(i, path), ObjectDiff)
											}
										}
									}
								}
								if i > len(valMod)-1 && !(parsing.MatchAny(valOrig[i], valMod)) {

									removed := parsing.RemovedDifference{Path: parsing.PathFormatter(parsing.PathSlice(i, path)),
										Value:                             valOrig[i]}
									ObjectDiff.Removed = append(ObjectDiff.Removed, removed)
								}
							}

						} else {
							for i := range valMod {
							Orig:
								for ii := range valOrig {
									if i == ii && reflect.DeepEqual(valOrig[ii], valMod[i]) {
										break Orig
									} else {
										if i != ii && reflect.DeepEqual(valOrig[ii], valMod[i]) {
											indexed := parsing.IndexDifference{OldIndex: ii, NewIndex: i, Value: valOrig[ii],
												Path: parsing.PathFormatter(path)}
											ObjectDiff.Indexes = append(ObjectDiff.Indexes, indexed)
											break Orig
										} else if i == ii && !(reflect.DeepEqual(valOrig[ii], valMod[i])) {
											if reflect.TypeOf(valOrig[ii]).Kind() == reflect.String || reflect.TypeOf(valMod[i]).Kind() == reflect.String {
												changed := parsing.ChangedDifference{Path: parsing.PathFormatter(parsing.PathSlice(i, path)),
													OldValue: valOrig[i], NewValue: valMod[i]}
												ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
												break Orig
											} else {
												ObjectDiff = recursion(parsing.Remarshal(valOrig[i]),
													parsing.Remarshal(valMod[i]), parsing.PathSlice(i, path), ObjectDiff)
											}
										}
									}
								}
								if i > len(valOrig)-1 && !(parsing.MatchAny(valMod[i], valOrig)) {

									added := parsing.AddedDifference{Path: parsing.PathFormatter(parsing.PathSlice(i, path)),
										Value:                             valMod[i]}
									ObjectDiff.Added = append(ObjectDiff.Added, added)
								}
							}
						}
					} else {
						// If both slice lengths are equal
						for i := range valOrig {
							if !(reflect.DeepEqual(valOrig[i], valMod[i])) {
								if reflect.TypeOf(valOrig[i]).Kind() == reflect.String || reflect.TypeOf(valMod[i]).Kind() == reflect.String {

									changed := parsing.ChangedDifference{Path: parsing.PathFormatter(parsing.PathSlice(i, path)),
										OldValue: valOrig[i], NewValue: valMod[i]}
									ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
								} else {

									ObjectDiff = recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]),
										parsing.PathSlice(i, path), ObjectDiff)
								}
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
	var ObjectDiff parsing.ConsumableDifference
	return recursion(original, modified, path, ObjectDiff)
}
