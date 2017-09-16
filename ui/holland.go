package ui

import (
	"encoding/json"
	"fmt"
	"github.com/beard1ess/gauss/operator"
	"github.com/beard1ess/gauss/parsing"
	"io"
	"log"
	"os"
	"reflect"
)

func check(action string, e error) {
	if e != nil {
		log.Fatal(action+" ", e)
	}
}

func Diff(

	origin string,
	modified string,
	output string,
	writer io.Writer,

) error {

	var json_original, json_modified parsing.Gaussian
	var path []string
	var objectDiff parsing.ConsumableDifference

	/* TODO WE WANT TO DO ALL OUR INIT STUFF IN THIS AREA */

	json_original.Read(origin)

	json_modified.Read(modified)

	if reflect.DeepEqual(json_original, json_modified) {
		writer.Write([]byte("No differences!"))
		os.Exit(0)
	} else {
		objectDiff = operator.Recursion(json_original.Data, json_modified.Data, path)
	}

	switch output {

	case "human":
		//writer.Write(format(objectDiff))

	case "machine":
		output, err := json.Marshal(objectDiff)
		check("sorry. ", err)

		writer.Write(output)

	default:
		fmt.Println("Output type unknown.")
		os.Exit(1)
	}

	return nil
}
