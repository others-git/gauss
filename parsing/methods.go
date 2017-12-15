package parsing

import (
	"io/ioutil"
	"log"
	"errors"

	"bytes"
	"encoding/json"
	"github.com/dimchansky/utfbom"

	"github.com/beard1ess/yaml"
	"hash/fnv"
	"io"
	"sort"
	"reflect"
	"fmt"
)


/*

GENERAL FUNCTIONS

 */

func check(action string, e error) {
	if e != nil {
		log.Fatal(action+" ", e)
	}
}

func forceSetter(input interface{}) (string, error) {
	if reflect.TypeOf(input).Kind() == reflect.Map {
		out,_ := json.Marshal(input)
		return string(out), nil
	}

	s, ok := input.(string)
	if !ok {
		s := fmt.Sprintf("unable to parse %v of %T as string", input, input)
		err := errors.New(s)
		return "", err
	}

	return s, nil
}

func hash(b []byte) uint32 {
	h := fnv.New32a()
	h.Write(b)
	return h.Sum32()
}

/*
CONSUMABLEDIFFERENCE TYPE FUNCTIONS
 */


func (c *ConsumableDifference) ReadFile(file string) error {

	// because go json refuses to deal with bom we need to strip it out
	f, err := ioutil.ReadFile(file)
	check(file, err)

	o, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(f)))
	check("Error encountered while trying to skip BOM: ", err)

	if err := json.Unmarshal(o, &c); err != nil {
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

// Sort each key in our object so marshaled object is also consistent
func (c *ConsumableDifference) Sort() error {

	// create 'sortable' string be combining fields that will always be present
	for i := range c.Changed {
		var buffer bytes.Buffer
		buffer.WriteString(c.Changed[i].Path)
		nv, err := forceSetter(c.Changed[i].NewValue)
		if err != nil {
			return err
		}
		buffer.WriteString(nv)
		c.Changed[i].sort = hash(buffer.Bytes())
	}
	for i := range c.Added {
		var buffer bytes.Buffer
		buffer.WriteString(c.Added[i].Path)
		av,err:= forceSetter(c.Added[i].Value)
		if err != nil {
			return err
		}
		buffer.WriteString(av)
		c.Added[i].sort = hash(buffer.Bytes())

	}
	for i := range c.Removed {
		var buffer bytes.Buffer
		buffer.WriteString(c.Removed[i].Path)
		rv, err := forceSetter(c.Removed[i].Value)
		if err != nil {
			return err
		}
		buffer.WriteString(rv)
		c.Removed[i].sort = hash(buffer.Bytes())

	}
	for i := range c.Indexes {
		var buffer bytes.Buffer
		buffer.WriteString(c.Indexes[i].Path)
		iv, err := forceSetter(c.Indexes[i].Value)
		if err != nil {
			return err
		}
		buffer.WriteString(iv)
		buffer.WriteString(string(c.Indexes[i].NewIndex))
		buffer.WriteString(string(c.Indexes[i].OldIndex))
		c.Indexes[i].sort = hash(buffer.Bytes())
	}
	sort.SliceStable(c.Changed, func(i, j int) bool { return c.Changed[i].sort < c.Changed[j].sort })
	sort.SliceStable(c.Added, func(i, j int) bool { return c.Added[i].sort < c.Added[j].sort })
	sort.SliceStable(c.Removed, func(i, j int) bool { return c.Removed[i].sort < c.Removed[j].sort })
	sort.SliceStable(c.Indexes, func(i, j int) bool { return c.Indexes[i].sort < c.Indexes[j].sort })
	return nil
}


// MarshalJSON Order and sort difference output for testing consistency
func (c *ConsumableDifference) MarshalJSON(input ...ConsumableDifference) ([]byte, error) {
	if input != nil {
		return json.Marshal(input)
	} else {
		//Since we don't actually care about the ordering of these, and they are slices, order by path to preserve tests
		c.Sort()
		return json.Marshal(c)
	}
}

/*

GAUSSIAN TYPE METHODS

 */

 // Read gaussian type reader method
func (g *Gaussian) Read(file string) error {
	var kvStore KeyValue
	// because go json refuses to deal with bom we need to strip it out
	f, err := ioutil.ReadFile(file)
	check(file, err)

	o, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(f)))
	check("Error encountered while trying to skip BOM: ", err)

	// We try to determine if json or yaml based on error :/
	err = json.Unmarshal(o, &kvStore)
	if err == nil {
		g.Data = kvStore
		g.Type = "JSON"
	} else {
		err = yaml.Unmarshal(o, &kvStore)
		if err == nil {
			g.Data = kvStore
			g.Type = "YAML"
		} else {
			err := errors.New("unable to parse file, confirm if valid JSON/YAML")
			return err
		}
	}
	return nil
}

// Write gaussian type writer method
func (g *Gaussian) Write(output io.Writer) error {

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
		err := errors.New("issue marshalling json/yaml to writer")
		return err
	}
	return nil
}
