package parsing

type Keyvalue map[string]interface{}
type Keyslice map[string][]Keyvalue

type Gaussian struct {
	Data Keyvalue // What we read into the struct
	Type string   // Json/Yaml

}

type RemovedDifference struct {
	Key   string `json:",omitempty"`
	Path  string
	Value interface{}
	sort  string `json:"-"`
}

type AddedDifference struct {
	Key   string `json:",omitempty"`
	Path  string
	Value interface{}
	sort  string `json:"-"`
}

type ChangedDifference struct {
	Key      string `json:",omitempty"`
	Path     string
	NewValue interface{}
	OldValue interface{}
	sort     string `json:"-"`
}

type IndexDifference struct {
	NewIndex int
	OldIndex int
	Path     string
	Value	 interface{}
	sort     string `json:"-"`
}

type ConsumableDifference struct {
	Changed []ChangedDifference `json:",omitempty"`
	Added   []AddedDifference   `json:",omitempty"`
	Removed []RemovedDifference `json:",omitempty"`
	Indexes []IndexDifference `json:",omitempty"`
}
