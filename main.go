package main

import (
	"encoding/json"
	"fmt"
	"github.com/beard1ess/gauss/operator"
	"github.com/beard1ess/gauss/parsing"
	"github.com/urfave/cli"
	"io"
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
	var patch, object string

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
					EnvVar: "ORIGINAL_OBJECT",
				},
				cli.StringFlag{
					Name: "modified, m",
					Usage: "Modified `OBJECT` to compare against",
					Value: "",
					EnvVar: "MODIFIED_OBJECT",
				},
				cli.StringFlag{
					Name: "output",
					Usage: "Output types available: human, machine",
					Value: "machine",
					EnvVar: "DIFF_OUTPUT",
				},
			},
			Action:  func(c *cli.Context) error {

				if c.String("origin") == "" {
					fmt.Print("ORIGIN is required!\n\n")
					cli.ShowCommandHelp(c, "diff")
					os.Exit(1)
				}

				if c.String("modified") == "" {
					fmt.Print("MODIFIED is required!\n\n")
					cli.ShowCommandHelp(c, "diff")
					os.Exit(1)
				}

				return diff(
					c.String("origin"),
					c.String("modified"),
					c.String("output"),
					os.Stdout,
				)
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
func diff(

	origin string,
	modified string,
	output string,
	writer io.Writer,

) error {

	var json_original, json_modified parsing.Keyvalue
	var path []string
	var objectDiff parsing.ConsumableDifference

	/* TODO WE WANT TO DO ALL OUR INIT STUFF IN THIS AREA */

	read, err := ioutil.ReadFile(origin)
	check(err)

	err = json.Unmarshal([]byte(read), &json_original)
	check(err)

	read, err = ioutil.ReadFile(modified)
	check(err)

	err = json.Unmarshal([]byte(read), &json_modified)
	check(err)

	if reflect.DeepEqual(json_original, json_modified) {
		fmt.Println("No differences!")
		os.Exit(0)
	} else {
		objectDiff = operator.Recursion(json_original, json_modified, path)
	}

	switch output {

	case "human":
		//writer.Write(format(objectDiff))

	case "machine":
		output, err := json.Marshal(objectDiff)
		check(err)

		writer.Write(output)

	default:
		fmt.Println("Output type unknown.")
		os.Exit(1)
	}

	return nil
}
