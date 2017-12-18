package operator

import (
	"github.com/beard1ess/gauss/parsing"
	"reflect"
)

func recursion(

	original parsing.KeyValue,
	modified parsing.KeyValue,
	path []string,
	ObjectDiff *parsing.ConsumableDifference,

)  {

	if reflect.DeepEqual(original, modified) {
		return
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

		recursion(original, modified, path, ObjectDiff)
		return

	} else if len(parsing.Slicer(original)) > 1 || len(parsing.Slicer(modified)) > 1 {

		for k := range original {
			 recursion(parsing.KeyValue{k: original[k]}, parsing.KeyValue{k: modified[k]}, path, ObjectDiff)
		}
		return
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
					return
					// Map handler
				} else if reflect.TypeOf(valOrig).Kind() == reflect.Map {
					// Update the working path
					path = append(path, k)
					recursion(parsing.Remarshal(valOrig), parsing.Remarshal(valMod), path, ObjectDiff)
					return
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

											if reflect.TypeOf(valOrig[i]).Kind() == reflect.String ||
												reflect.TypeOf(valMod[ii]).Kind() == reflect.String ||
												!(parsing.DoMapArrayKeysMatch(valOrig[i], valMod[ii])) {

												changed := parsing.ChangedDifference{Path: parsing.PathFormatter(parsing.SliceIndex(i, path)),
													OldValue: valOrig[i], NewValue: valMod[i]}
												ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
												break Mod

											} else {

												recursion(parsing.Remarshal(valOrig[i]),
													parsing.Remarshal(valMod[i]), parsing.SliceIndex(i, path), ObjectDiff)
											}
										}
									}
								}
								if i > len(valMod)-1 && !(parsing.MatchAny(valOrig[i], valMod)) {

									removed := parsing.RemovedDifference{Path: parsing.PathFormatter(parsing.SliceIndex(i, path)),
										Value: valOrig[i]}
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
												changed := parsing.ChangedDifference{Path: parsing.PathFormatter(parsing.SliceIndex(i, path)),
													OldValue: valOrig[i], NewValue: valMod[i]}
												ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
												break Orig
											} else {

												recursion(parsing.Remarshal(valOrig[i]),
													parsing.Remarshal(valMod[i]), parsing.SliceIndex(i, path), ObjectDiff)
											}
										}
									}
								}
								if i > len(valOrig)-1 && !(parsing.MatchAny(valMod[i], valOrig)) {

									added := parsing.AddedDifference{Path: parsing.PathFormatter(parsing.SliceIndex(i, path)),
										Value: valMod[i]}
									ObjectDiff.Added = append(ObjectDiff.Added, added)
								}
							}
						}
					} else {
						// If both slice lengths are equal
						for i := range valOrig {
							if !(reflect.DeepEqual(valOrig[i], valMod[i])) {
								if reflect.TypeOf(valOrig[i]).Kind() == reflect.String || reflect.TypeOf(valMod[i]).Kind() == reflect.String {

									changed := parsing.ChangedDifference{Path: parsing.PathFormatter(parsing.SliceIndex(i, path)),
										OldValue: valOrig[i], NewValue: valMod[i]}
									ObjectDiff.Changed = append(ObjectDiff.Changed, changed)
								} else if reflect.TypeOf(valOrig[i]).Kind() == reflect.Slice || reflect.TypeOf(valMod[i]).Kind() == reflect.Slice {

									changed := parsing.ChangedDifference{Path: parsing.PathFormatter(parsing.PathSlice(i, path)),
										OldValue: valOrig[i], NewValue: valMod[i]}
									ObjectDiff.Changed = append(ObjectDiff.Changed, changed)

								} else {

									recursion(parsing.Remarshal(valOrig[i]), parsing.Remarshal(valMod[i]),
										parsing.SliceIndex(i, path), ObjectDiff)
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
		return
	}
}

// Recursion wrapper for primary recursion function to find differences
func Recursion(original parsing.KeyValue, modified parsing.KeyValue, path []string) parsing.ConsumableDifference {
	var ObjectDiff parsing.ConsumableDifference
	recursion(original, modified, path, &ObjectDiff)
	return ObjectDiff
}
