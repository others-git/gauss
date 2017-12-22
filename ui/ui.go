package ui

import (
	"github.com/beard1ess/gauss/operator"
	"github.com/beard1ess/gauss/parsing"
	"io"
	"reflect"
	"errors"
	"github.com/jmespath/go-jmespath"
	"fmt"
	"encoding/json"
)

/*
ui package is for all interfacing and commands we expose
*/

// Diff handle file inputs and pass to function to find differences
func Diff(

	origin string,
	modified string,
	output string,
	inputDiffPath string,
	writer io.Writer,

) error {

	var jsonOriginal, jsonModified parsing.Gaussian
	var path []string
	var objectDiff *parsing.ConsumableDifference
	var err error


	/* TODO WE WANT TO DO ALL OUR INIT STUFF IN THIS AREA */

	if err := jsonOriginal.Read(origin) ; err != nil {
		return err
	}

	if err := jsonModified.Read(modified) ; err != nil {
		return err
	}


	// Validate jmespath expression and move into path if exists
	if len(inputDiffPath) > 0 {
		_, err :=  jmespath.Compile(inputDiffPath)
		if err != nil {
			nErr := fmt.Errorf("failed to compile provided path: %T", err)
			return nErr
		}
		jsonOriginal.Data,err = jmespath.Search(inputDiffPath, jsonOriginal.Data)
		if jsonOriginal.Data == nil {
			err := errors.New("difference path returned nil object")
			return err
		} else if err != nil {
			nErr := fmt.Errorf("error pathing to object in original: %T", err)
			return nErr
		}
		jsonModified.Data,err = jmespath.Search(inputDiffPath, jsonModified.Data)
		if jsonModified.Data == nil {
			err := errors.New("difference path returned nil object")
			return err
		} else if err != nil {
			nErr := fmt.Errorf("error pathing to object in modified: %T", err)
			return nErr
		}
	}

	if reflect.DeepEqual(jsonOriginal, jsonModified) {
		writer.Write([]byte("No differences!"))
		return nil
	} else {
		objectDiff, err = operator.Recursion(
			jsonOriginal.Data.(map[string]interface{}),
			jsonModified.Data.(map[string]interface{}),
			path)
		if err != nil {
			return err
		}
	}


	switch output {

	case "formatted":
		//writer.Write(format(objectDiff))

	case "raw":
		objectDiff.Sort()
		output, err := json.Marshal(objectDiff)
		if err != nil {
			return err
		}

		writer.Write(output)

	default:
		err := fmt.Errorf("output type unknown: %T", output)
		return err
	}

	return nil
}

// Patch unused
func Patch(

	patch string,
	original string,
	output string,
	skipKeys string,
	writer io.Writer,

) error {
	var patcher parsing.ConsumableDifference
	var originObject parsing.Gaussian

	patcher.Read(patch)
	//parsing.Format(patcher)

	originObject.Read(original)

	operator.Patch(&patcher, &originObject)

	/*
	switch output {

	case "formatted":
		//writer.Write(format(objectDiff))

	case "raw":

		objectDiff.Sort()
		output, err := json.Marshal(objectDiff)
		if err != nil {
			return err
		}

		writer.Write(output)

	default:
		err := fmt.Errorf("output type unknown: %T", output)
		return err
	}
	*/

	return nil
}


























