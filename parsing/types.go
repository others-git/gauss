package parsing


// KeyValue map of string to interface
type KeyValue map[string]interface{}

// KeySlice map of string to slice of KeyValue
type KeySlice map[string][]KeyValue

// Gaussian wrapper to handle inputs
type Gaussian struct {
	Data interface{} // What we read into the struct
	Type string   // Json/Yaml
}

// RemovedDifference sub struct for removed objects
type RemovedDifference struct {
	Key   string `json:",omitempty"`
	Path  string
	Value interface{}
	sort  uint32
}

// AddedDifference sub struct for added objects
type AddedDifference struct {
	Key   string `json:",omitempty"`
	Path  string
	Value interface{}
	sort  uint32
}

// ChangedDifference sub struct for changed objects
type ChangedDifference struct {
	Key      string `json:",omitempty"`
	Path     string
	NewValue interface{}
	OldValue interface{}
	sort     uint32
}

// IndexDifference sub struct for list index changes
type IndexDifference struct {
	NewIndex int
	OldIndex int
	Path     string
	Value	 interface{}
	sort     uint32
}

// ConsumableDifference eventual return object
type ConsumableDifference struct {
	Changed []ChangedDifference `json:",omitempty"`
	Added   []AddedDifference   `json:",omitempty"`
	Removed []RemovedDifference `json:",omitempty"`
	Indexes []IndexDifference `json:",omitempty"`
}
