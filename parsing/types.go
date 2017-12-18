package parsing


// KeyValue map of string to interface
type KeyValue map[string]interface{}

// KeySlice map of string to slice of KeyValue
type KeySlice map[string][]KeyValue

// Gaussian reader type to handle reading the file and determine if json/yaml
type Gaussian struct {
	Data KeyValue // What we read into the struct
	Type string   // Json/Yaml

}

// RemovedDifference removed items from the object
type RemovedDifference struct {
	Key   string `json:",omitempty"`
	Path  string
	Value interface{}
	sort  uint32 `json:"-"`
}

// AddedDifference added items from the object
type AddedDifference struct {
	Key   string `json:",omitempty"`
	Path  string
	Value interface{}
	sort  uint32 `json:"-"`
}

// ChangedDifference changed items from the object
type ChangedDifference struct {
	Key      string `json:",omitempty"`
	Path     string
	NewValue interface{}
	OldValue interface{}
	sort     uint32 `json:"-"`
}

// IndexDifference try to determine when array items only change index
type IndexDifference struct {
	NewIndex int
	OldIndex int
	Path     string
	Value	 interface{}
	sort     uint32 `json:"-"`
}

// ConsumableDifference eventual return object
type ConsumableDifference struct {
	Changed []ChangedDifference `json:",omitempty"`
	Added   []AddedDifference   `json:",omitempty"`
	Removed []RemovedDifference `json:",omitempty"`
	Indexes []IndexDifference `json:",omitempty"`
}
