package main

import (
	"fmt"
	"github.com/beard1ess/gauss/ui"
	"github.com/urfave/cli"
	"os"
)

func main() {

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
					Name:   "origin, o",
					Usage:  "Original `OBJECT` to compare against",
					Value:  "",
					EnvVar: "ORIGINAL_OBJECT",
				},
				cli.StringFlag{
					Name:   "modified, m",
					Usage:  "Modified `OBJECT` to compare against",
					Value:  "",
					EnvVar: "MODIFIED_OBJECT",
				},
				cli.StringFlag{
					Name:   "output",
					Usage:  "Output types available: human, machine",
					Value:  "machine",
					EnvVar: "DIFF_OUTPUT",
				},
				cli.StringFlag{
					Name:  "in, i",
					Usage: "Because some roads you shouldn't go down. Because maps used to say, \"There be dragons here.\" Now they don't. But that don't mean the dragons aren't there.",
					Value: "",
				},
			},
			Action: func(c *cli.Context) error {

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

				return ui.Diff(
					c.String("origin"),
					c.String("modified"),
					c.String("output"),
					os.Stdout,
				)

			},
		},
		/* NOT IMPLEMENTED
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
		*/
	}

	app.Run(os.Args)

}
