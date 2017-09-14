package parsing

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dimchansky/utfbom"
	"bytes"
	"encoding/json"

	"github.com/beard1ess/yaml"
	"os"
	"io"
)

func check(action string, e error) {
	if e != nil {
		log.Fatal(action + " ", e)
	}
}

type Gaussian struct {

	Data Keyvalue // What we read into the struct
	Type string // Json/Yaml

}

func (g *Gaussian) Read(input string) {
	var kv_store Keyvalue
	// because go json refuses to deal with bom we need to strip it out
	f, err := ioutil.ReadFile(input)
	check(input, err)

	o,err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(f)))
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

// I wrote this and realized it may not be useful
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
		fmt.Println("Someout TYPE is messed up for Gaussian struct.")
		os.Exit(9001)
	}
}
