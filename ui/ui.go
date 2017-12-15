package ui

import (
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

	var jsonOriginal, jsonModified parsing.Gaussian
	var path []string
	var objectDiff parsing.ConsumableDifference

	/* TODO WE WANT TO DO ALL OUR INIT STUFF IN THIS AREA */

	if err := jsonOriginal.Read(origin) ; err != nil {
		return err
	}

	jsonModified.Read(modified)

	if reflect.DeepEqual(jsonOriginal, jsonModified) {
		writer.Write([]byte("No differences!"))
		os.Exit(0)
	} else {
		objectDiff = operator.Recursion(jsonOriginal.Data, jsonModified.Data, path)
	}

	switch output {

	case "human":
		//writer.Write(format(objectDiff))

	case "machine":
		output, err := objectDiff.MarshalJSON()

		check("sorry. ", err)

		writer.Write(output)

	default:
		fmt.Println("Output type unknown.")
		os.Exit(1)
	}

	return nil
}


func Patch(

	patch string,
	origin string,
	output string,
	writer io.Writer,


) error {
	var patcher parsing.ConsumableDifference
	patcher.ReadFile(patch)

	fmt.Println(patcher.Added[0].Value)

	return nil
}