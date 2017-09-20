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
	"reflect"
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

func forcesertter(input interface{}) string {
	if reflect.TypeOf(input).Kind() == reflect.Map {
		out,_ := json.Marshal(input)
		return string(out)
	}
	return input.(string)
}

func (c *ConsumableDifference) Sort() {

	// create 'sortable' string be combining fields that will always be present
	for i := range c.Changed {
		var buffer bytes.Buffer
		buffer.WriteString(c.Changed[i].Path)
		buffer.WriteString(forcesertter(c.Changed[i].NewValue))
		buffer.WriteString(forcesertter(c.Changed[i].OldValue))
		c.Changed[i].sort = buffer.String()
	}
	for i := range c.Added {
		var buffer bytes.Buffer
		buffer.WriteString(c.Added[i].Path)
		buffer.WriteString(forcesertter(c.Added[i].Value))
		c.Added[i].sort = buffer.String()
	}
	for i := range c.Removed {
		var buffer bytes.Buffer
		buffer.WriteString(c.Removed[i].Path)
		buffer.WriteString(forcesertter(c.Removed[i].Value))
		c.Removed[i].sort = buffer.String()
	}
	for i := range c.Indexes {
		var buffer bytes.Buffer
		buffer.WriteString(c.Indexes[i].Path)
		buffer.WriteString(forcesertter(c.Indexes[i].Value))
		buffer.WriteString(string(c.Indexes[i].NewIndex))
		buffer.WriteString(string(c.Indexes[i].OldIndex))
		c.Indexes[i].sort = buffer.String()
	}
	sort.SliceStable(c.Changed, func(i, j int) bool { return c.Changed[i].sort < c.Changed[j].sort })
	sort.SliceStable(c.Added, func(i, j int) bool { return c.Added[i].sort < c.Added[j].sort })
	sort.SliceStable(c.Removed, func(i, j int) bool { return c.Removed[i].sort < c.Removed[j].sort })
	sort.SliceStable(c.Indexes, func(i, j int) bool { return c.Indexes[i].sort < c.Indexes[j].sort })
}

func (c *ConsumableDifference) MarshalJSON(input ...ConsumableDifference) ([]byte, error) {
	if input != nil {
		return json.Marshal(input)
	} else {
		//Since we don't actually care about the ordering of these, and they are slices, order by path to preserve tests
		c.Sort()
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
