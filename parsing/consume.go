package parsing

import (
	"fmt"
	"io/ioutil"
	"log"

	"bytes"
	"encoding/json"
	"github.com/dimchansky/utfbom"

	"github.com/beard1ess/yaml"
	"io"
	"os"
	"sort"
)

type Keyvalue map[string]interface{}
type Keyslice map[string][]Keyvalue

func check(action string, e error) {
	if e != nil {
		log.Fatal(action+" ", e)
	}
}

type RemovedDifference struct {
	Key   string `json:",omitempty"`
	Path  string
	Value interface{}
}

type AddedDifference struct {
	Key   string `json:",omitempty"`
	Path  string
	Value interface{}
}

type ChangedDifference struct {
	Key      string `json:",omitempty"`
	Path     string
	NewValue interface{}
	OldValue interface{}
}

type IndexDifference struct {
	NewIndex int
	OldIndex int
	Path     string
	Value	 interface{}
}

type ConsumableDifference struct {
	Changed []ChangedDifference `json:",omitempty"`
	Added   []AddedDifference   `json:",omitempty"`
	Removed []RemovedDifference `json:",omitempty"`
	Indexes []IndexDifference `json:",omitempty"`
}

func (c *ConsumableDifference) ReadFile(file string) error {

	// because go json refuses to deal with bom we need to strip it out
	f, err := ioutil.ReadFile(file)
	check(file, err)

	o, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(f)))
	check("Error encountered while trying to skip BOM: ", err)

	if err := json.Unmarshal(o, &c) ; err != nil {
		return err
	}

	return nil
}

/* UNUSED, MAYBE NOT USEFUL AT ALL, WILL COME BACK TO LATER.
 * PROBABLY NEED THIS TO GIVE INTERFACE TO THE STRUCT FOR PROGRAMS
func (c *ConsumableDifference) UnmarshalJSON(input ...interface{}) error {
	if input == nil {

	} else {

	}

	return nil
}
*/

func (c *ConsumableDifference) MarshalJSON(input ...ConsumableDifference) ([]byte, error) {
	if input != nil {
		return json.Marshal(input)
	} else {
		//Since we don't actually care about the ordering of these, and they are slices, order by path to preserve tests
		sort.SliceStable(c.Changed, func(i, j int) bool { return c.Changed[i].Path < c.Changed[j].Path })
		sort.SliceStable(c.Added, func(i, j int) bool { return c.Added[i].Path < c.Added[j].Path })
		sort.SliceStable(c.Removed, func(i, j int) bool { return c.Removed[i].Path < c.Removed[j].Path })
		sort.SliceStable(c.Indexes, func(i, j int) bool { return c.Indexes[i].Path < c.Indexes[j].Path })
		return json.Marshal(c)
	}

	return nil, nil
}

type Gaussian struct {
	Data Keyvalue // What we read into the struct
	Type string   // Json/Yaml

}

func (g *Gaussian) Read(file string) {
	var kv_store Keyvalue
	// because go json refuses to deal with bom we need to strip it out
	f, err := ioutil.ReadFile(file)
	check(file, err)

	o, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(f)))
	check("Error encountered while trying to skip BOM: ", err)

	// We try to determine if json or yaml based on error :/
	err = json.Unmarshal(o, &kv_store)
	if err == nil {
		g.Data = kv_store
		g.Type = "JSON"
	} else {
		err = yaml.Unmarshal(o, &kv_store)
		if err == nil {
			g.Data = kv_store
			g.Type = "YAML"
		} else {
			fmt.Println("Unparseable file type presented")
			os.Exit(2)
		}
	}
}

// I wrote this and realized it may not be useful, pass a writer to the function and it will marshal and write out the data
func (g *Gaussian) Write(output io.Writer) {

	switch g.Type {
	case "JSON":

		o, err := json.Marshal(g.Data)
		check("Gaussian marshal error. ", err)
		output.Write(o)

	case "YAML":

		o, err := yaml.Marshal(g.Data)
		check("Gaussian marshal error. ", err)
		output.Write(o)

	default:
		fmt.Println("Somehow TYPE is messed up for Gaussian struct.")
		os.Exit(9001)
	}
}
