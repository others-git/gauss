package main

import (
	"os"
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"log"
	"reflect"
	"regexp"
	"github.com/beard1ess/gauss/parsing"
	"github.com/beard1ess/gauss/operator"

)

var (
	FormattedDiff parsing.Keyslice

)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}



func format(input parsing.ConsumableDifference) parsing.Keyvalue {
	var return_value parsing.Keyvalue

	FormattedDiff = nil
	/*
	for i := range input["Changed"] {
		path_builder(input["Changed"][i]["Path"].([]string))
	}
	for i := range input["Added"] {
		path_builder(input["Added"][i]["Path"].([]string))
	}
	for i := range input["Removed"] {
		path_builder(input["Removed"][i]["Path"].([]string))

	}
	*/

	return return_value
}

func path_builder(path []string)  parsing.Keyvalue{
	var object parsing.Keyvalue
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
	fmt.Println(path)
	return object
}

func main() {
	var patch, object, original_obj, modified_obj string

	app := cli.NewApp()
	app.Name = "Gauss"
	app.Version = "0.1"
	app.Usage = "Objected-based difference and patching tool."

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
					Value: "machine",
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
				var json_original, json_modified parsing.Keyvalue
				var path []string
				var ObjectDiff parsing.ConsumableDifference
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
					ObjectDiff = operator.Recursion(json_original, json_modified, path)
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


