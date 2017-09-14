package main

import (
	"encoding/json"
	"fmt"
	"github.com/beard1ess/gauss/operator"
	"github.com/beard1ess/gauss/parsing"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	var patch, object, original_obj, modified_obj string

	app := cli.NewApp()
	app.Name = "Gauss"
	app.Version = "0.1"
	app.Usage = "Objected-based difference and patching tool."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "test, t",
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
					Name:        "origin, o",
					Usage:       "Original `OBJECT` to compare against",
					Value:       "",
					Destination: &original_obj,
					EnvVar:      "ORIGINAL_OBJECT",
				},
				cli.StringFlag{
					Name:        "modified, m",
					Usage:       "Modified `OBJECT` to compare against",
					Value:       "",
					Destination: &modified_obj,
					EnvVar:      "MODIFIED_OBJECT",
				},
				cli.StringFlag{
					Name:   "output",
					Usage:  "Output types available: human, machine",
					Value:  "machine",
					EnvVar: "DIFF_OUTPUT",
				},
			},
			Action: func(c *cli.Context) error {
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

				read, err := ioutil.ReadFile(original_obj)
				check(err)
				err = json.Unmarshal([]byte(read), &json_original)
				check(err)
				read, err = ioutil.ReadFile(modified_obj)
				check(err)
				err = json.Unmarshal([]byte(read), &json_modified)
				check(err)
				if reflect.DeepEqual(json_original, json_modified) {
					fmt.Println("No differences!")
					os.Exit(0)
				} else {
					ObjectDiff = operator.Recursion(json_original, json_modified, path)
				}
				if c.String("output") == "human" {
					parsing.Format(ObjectDiff)
				} else if c.String("output") == "machine" {
					output, _ := json.Marshal(ObjectDiff)
					os.Stdout.Write(output)
				} else {
					fmt.Println("Output type unknown.")
					os.Exit(1)
				}

				return nil
			},
		},
		{
			Name:    "patch",
			Aliases: []string{"p"},
			Usage:   "Apply patch file to json object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "patch, p",
					Usage:       "`PATCH` the OBJECT",
					Value:       "",
					Destination: &patch,
				},
				cli.StringFlag{
					Name:        "object, o",
					Usage:       "`OBJECT` to PATCH",
					Value:       "",
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
