package main

import (
	"os"
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"regexp"
)

var (
	//ObjectDiff Keyslice
	ObjectDiff = make(Keyslice)
	FormattedDiff = make(Keyvalue)
)

type Keyvalue map[string]interface{}
type Keyslice map[string][]Keyvalue

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Remarshal(input interface{}) Keyvalue {
	var back Keyvalue
	out,_ := json.Marshal(input)
	_ = json.Unmarshal([]byte(out), &back)
	return back
}

func ListStripper(input Keyvalue) []string {
	var r []string
	for key := range input {
		r = append(r, key)
	}
	return r
}

func IndexOf(inputList []string, inputKey string) int {
	for i, v := range inputList {
		if v == inputKey {
			return i
		}
	}
	return -1
}



func recursion(original Keyvalue, modified Keyvalue, path []string) {

	kListModified := ListStripper(modified)
	kListOriginal := ListStripper(original)

	if len(kListModified) > 1 || len(kListOriginal) > 1 {
		proc := true
		for k, v := range original {
			if IndexOf(kListModified, k) == -1 {
				ObjectDiff["Removed"] = append(ObjectDiff["Removed"],Keyvalue{"Path": path, "Key": k, "Value":v})
				proc = false
			}
		}

		for k, v := range modified {
			if IndexOf(kListOriginal, k) == -1 {
				ObjectDiff["Added"] = append(ObjectDiff["Added"],Keyvalue{"Path": path, "Key": k, "Value":v})
				proc = false
			}
		}
		if proc {
			for k := range original {
				recursion(Keyvalue{k:original[k]},Keyvalue{k:modified[k]},path)
			}
		}
		return
	}
	for k := range original {
		var valOrig, valMod interface{}

		if reflect.TypeOf(original).Name() == "string" {

			valOrig = original
		} else {
			valOrig = original[k]
		}
		if reflect.TypeOf(modified).Name() == "string" {
			valMod = modified
		} else {
			valMod = modified[k]
		}

		if !(reflect.DeepEqual(valMod, valOrig)) {

			if reflect.TypeOf(valOrig).Kind() == reflect.Map {
				npath := append(path, k)
				//recursion(valOrig.(Keyvalue), valMod.(Keyvalue), npath)
				recursion(Remarshal(valOrig), Remarshal(valMod), npath)
				return
			} else if reflect.TypeOf(valOrig).Kind() == reflect.Slice {
				valOrig,_ := valOrig.([]interface{})
				valMod,_ := valMod.([]interface{})
				if len(valOrig) != len(valMod) {
					// do things
					fmt.Println("Unexpected spot, array lengths are different for value, not implemented lol")
					os.Exit(1)
				} else {
					for i := range valOrig {
						if !(reflect.DeepEqual(valMod[i], valOrig[i])) {
							npath := append(path, "{Index:"+strconv.Itoa(i)+"}")
							ObjectDiff["Changed"] = append(ObjectDiff["Changed"],Keyvalue{"Path": npath,
								"Key": k, "oldValue":valOrig[i],"newValue":valMod[i]})
							return
						}
					}
				}
			} else {
				ObjectDiff["Changed"] = append(ObjectDiff["Changed"],Keyvalue{"Path": path, "Key": k,
					"oldValue":valOrig,"newValue":valMod})
				return
			}
		}
		return
	}
	return
}


func format(input Keyslice) Keyvalue {
	var return_value Keyvalue




	for i := range input["Changed"] {
		path_builder(input["Changed"][i]["Path"].([]string))
	}
	for i := range input["Added"] {
		path_builder(input["Added"][i]["Path"].([]string))
	}
	for i := range input["Removed"] {
		path_builder(input["Removed"][i]["Path"].([]string))

	}


	return return_value
}

func path_builder(path []string)  Keyvalue{
	var object Keyvalue
	FormattedDiff = nil
	r, _ := regexp.Compile("[0-9]+")
	//path_length := len(path)
	for i:= range path {
		if ok,_ := regexp.MatchString("{Index:[0-9]+}", path[i]); ok {
			index := r.FindString(path[i])
			fmt.Println(index)
		} else {

		}
	}


	fmt.Println(path)
	return object
}

func main() {
	var patch, object, original_obj, modified_obj string

	app := cli.NewApp()
	app.Name = "JsonDiffer"
	app.Version = "0.1"
	app.Usage = "Used to get an object-based difference between two json objects."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "test, t",
			Usage: "just taking up space",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "diff",
			Aliases: []string{"d"},
			Usage:   "Diff json objects",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "origin, o",
					Usage: "Original `OBJECT` to compare against",
					Value: "",
					Destination: &original_obj,
					EnvVar: "ORIGINAL_OBJECT",
				},
				cli.StringFlag{
					Name: "modified, m",
					Usage: "Modified `OBJECT` to compare against",
					Value: "",
					Destination: &modified_obj,
					EnvVar: "MODIFIED_OBJECT",
				},
				cli.StringFlag{
					Name: "output",
					Usage: "Output types available: human, machine",
					Value: "human",
					EnvVar: "DIFF_OUTPUT",
				},
				/*
				cli.StringFlag{
					Name: "output, O",
					Usage: "File output location",
					Value: "",
					Destination: &modified_obj,
				},
				*/
			},
			Action:  func(c *cli.Context) error {
				var json_original, json_modified Keyvalue
				var path []string
				if original_obj == "" {
					fmt.Print("ORIGIN is required!\n\n")
					cli.ShowCommandHelp(c, "diff")
					os.Exit(1)
				}
				if modified_obj == "" {
					fmt.Print("MODIFIED is required!\n\n")
					cli.ShowCommandHelp(c, "diff")
					os.Exit(1)
				}

				/* TODO WE WANT TO DO ALL OUR INIT STUFF IN THIS AREA */

				/*
				ObjectDiff["Changed"] = []Keyvalue{}
				ObjectDiff["Added"] = []Keyvalue{}
				ObjectDiff["Removed"] = []Keyvalue{}
				*/

				read,err := ioutil.ReadFile(original_obj)
				check(err)
				_ = json.Unmarshal([]byte(read), &json_original)

				read,err = ioutil.ReadFile(modified_obj)
				check(err)
				_ = json.Unmarshal([]byte(read), &json_modified)


				if reflect.DeepEqual(json_original, json_modified) {
					fmt.Println("No differences!")
					os.Exit(0)
				} else {
					recursion(json_original, json_modified, path)
				}

				if c.String("output") == "human" {
					format(ObjectDiff)
				} else if c.String("output") == "machine" {
					output,_ := json.Marshal(ObjectDiff)
					os.Stdout.Write(output)
				} else {
					fmt.Println("Output type unknown.")
					os.Exit(1)
				}

				return nil
			},
		},
		{
			Name: "patch",
			Aliases: []string{"p"},
			Usage:	"Apply patch file to json object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "patch, p",
					Usage: "`PATCH` the OBJECT",
					Value: "",
					Destination: &patch,
				},
				cli.StringFlag{
					Name: "object, o",
					Usage: "`OBJECT` to PATCH",
					Value: "",
					Destination: &object,
				},
			},
			Action: func(c *cli.Context) error {

				return nil
			},
		},
	}

	app.Run(os.Args)

}


