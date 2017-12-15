package parsing

import (
	"io/ioutil"
	"log"
	"errors"

	"bytes"
	"encoding/json"
	"github.com/dimchansky/utfbom"

	"github.com/beard1ess/yaml"
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

func forcesertter(input interface{}) (string, error) {
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

/*
CONSUMABLEDIFFERENCE TYPE FUNCTIONS
 */

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

// Sort each key in our object so marshaled object is also consistent
func (c *ConsumableDifference) Sort() error {

	// create 'sortable' string be combining fields that will always be present
	for i := range c.Changed {
		var buffer bytes.Buffer
		buffer.WriteString(c.Changed[i].Path)
		fnv, err := forcesertter(c.Changed[i].NewValue)
		if err != nil {
			return err
		}
		buffer.WriteString(fnv)
		fov, err := forcesertter(c.Changed[i].OldValue)
		if err != nil {
			return err
		}
		buffer.WriteString(fov)
		c.Changed[i].sort = buffer.String()
	}
	for i := range c.Added {
		var buffer bytes.Buffer
		buffer.WriteString(c.Added[i].Path)
		fv, err := forcesertter(c.Added[i].Value)
		if err != nil {
			return err
		}
		buffer.WriteString(fv)
		c.Added[i].sort = buffer.String()
	}
	for i := range c.Removed {
		var buffer bytes.Buffer
		buffer.WriteString(c.Removed[i].Path)
		fv, err := forcesertter(c.Removed[i].Value)
		if err != nil {
			return err
		}
		buffer.WriteString(fv)
		c.Removed[i].sort = buffer.String()
	}
	for i := range c.Indexes {
		var buffer bytes.Buffer
		buffer.WriteString(c.Indexes[i].Path)
		fv, err := forcesertter(c.Indexes[i].Value)
		if err != nil {
			return err
		}
		buffer.WriteString(fv)
		buffer.WriteString(string(c.Indexes[i].NewIndex))
		buffer.WriteString(string(c.Indexes[i].OldIndex))
		c.Indexes[i].sort = buffer.String()
	}
	sort.SliceStable(c.Changed, func(i, j int) bool { return c.Changed[i].sort < c.Changed[j].sort })
	sort.SliceStable(c.Added, func(i, j int) bool { return c.Added[i].sort < c.Added[j].sort })
	sort.SliceStable(c.Removed, func(i, j int) bool { return c.Removed[i].sort < c.Removed[j].sort })
	sort.SliceStable(c.Indexes, func(i, j int) bool { return c.Indexes[i].sort < c.Indexes[j].sort })
	return nil
}

// Order and sort difference output for testing consistency
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

/*

GAUSSIAN TYPE METHODS

 */

func (g *Gaussian) Read(file string) error {
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
			error := errors.New("unable to parse file, confirm if valid JSON/YAML")
			return error
		}
	}
	return nil
}

// I wrote this and realized it may not be useful, pass a writer to the function and it will marshal and write out the data
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
		error := errors.New("issue marshalling json/yaml to writer")
		return error
	}
	return nil
}
