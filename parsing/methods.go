package parsing

import (
	"io/ioutil"
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
func forceSetter(input interface{}) (string, error) {
	if reflect.TypeOf(input).Kind() == reflect.Map {
		out,_ := json.Marshal(input)
		return string(out), nil
	}
	s, ok := input.(string)
	if !ok {
		err := errors.New(fmt.Sprintf("unable to parse %v of type %T as string", input, input))
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

// ConsumableDifference.Read helper to read file that handles bom and unmarshals for our patch
func (c *ConsumableDifference) Read(file string) error {

	// because go json refuses to deal with bom we need to strip it out
	f, err := ioutil.ReadFile(file)
	if err != nil {
		nErr := errors.New(fmt.Sprintf("error reading file %T: %T", file, err))
		return nErr
	}

	o, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(f)))
	if err != nil {
		nErr := errors.New(fmt.Sprintf("error encountered while trying to skip BOM: %T", err))
		return nErr
	}

	if err := json.Unmarshal(o, &c); err != nil {
		return err
	}

	return nil
}

// ConsumableDifference.Sort each key in our object so marshaled object is also consistent
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

/*

GAUSSIAN TYPE METHODS

 */

 // Gaussian.Read gaussian type reader method for raw json or yaml
func (g *Gaussian) Read(file string) error {
	var store interface{}
	// because go json refuses to deal with bom we need to strip it out
	f, err := ioutil.ReadFile(file)
	if err != nil {
		nErr := errors.New(fmt.Sprintf("error reading file %T: %T", file, err))
		return nErr
	}

	o, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(f)))
	if err != nil {
		nErr := errors.New(fmt.Sprintf("error encountered while trying to skip BOM: %T", err))
		return nErr
	}

	// We try to determine if json or yaml based on error :/
	err = json.Unmarshal(o, &store)
	if err == nil {
		g.Data = store
		g.Type = "JSON"
	} else {
		err = yaml.Unmarshal(o, &store)
		if err == nil {
			g.Data = store
			g.Type = "YAML"
		} else {
			err := errors.New("unable to parse file, confirm if valid JSON/YAML")
			return err
		}
	}
	return nil
}

// Gaussian.Write marshals data based on type and outputs to writer
func (g *Gaussian) Write(output io.Writer) error {

	switch g.Type {
	case "JSON":

		o, err := json.Marshal(g.Data)
		if err != nil {
			nErr := errors.New(fmt.Sprintf("error marshalling input: %T", err))
			return nErr
		}
		output.Write(o)

	case "YAML":

		o, err := yaml.Marshal(g.Data)
		if err != nil {
			nErr := errors.New(fmt.Sprintf("error marshalling input: %T", err))
			return nErr
		}
		output.Write(o)

	default:
		err := errors.New("issue marshalling json/yaml to writer")
		return err
	}
	return nil
}
