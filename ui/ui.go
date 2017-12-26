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
	"regexp"
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
	regSkipKeys string,
	writer io.Writer,

) error {

	var jsonOriginal, jsonModified parsing.Gaussian
	var path []string
	var objectDiff *parsing.ConsumableDifference
	var regSkip *regexp.Regexp
	var err error


	/* TODO WE WANT TO DO ALL OUR INIT STUFF IN THIS AREA */

	if err := jsonOriginal.ReadFile(origin) ; err != nil {
		return err
	}

	if err := jsonModified.ReadFile(modified) ; err != nil {
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
			err := fmt.Errorf("path %v returned nil for %v", inputDiffPath, origin)
			return err
		} else if err != nil {
			nErr := fmt.Errorf("error pathing to object in original: %T", err)
			return nErr
		}
		jsonModified.Data,err = jmespath.Search(inputDiffPath, jsonModified.Data)
		if jsonModified.Data == nil {
			err := fmt.Errorf("path %v returned nil for %v", inputDiffPath, modified)
			return err
		} else if err != nil {
			nErr := fmt.Errorf("error pathing to object in modified: %T", err)
			return nErr
		}
	}

	// compile provided regexp if provided
	if len(regSkipKeys) > 0 {
		regSkip,err = regexp.Compile(regSkipKeys)
		if err != nil {
			return err
		}
	}

	if reflect.DeepEqual(jsonOriginal, jsonModified) {
		writer.Write([]byte("No differences!"))
		return nil
	} else {
		objectDiff, err = operator.Recursion(
			jsonOriginal.Data.(map[string]interface{}),
			jsonModified.Data.(map[string]interface{}),
			path,
			regSkip,
			)
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
	regSkipKeys string,
	inputDiffPath string,
	writer io.Writer,

) error {
	var patcher parsing.ConsumableDifference
	var originObject parsing.Gaussian
	var regSkip *regexp.Regexp
	var err error

	patcher.ReadFile(patch)
	//parsing.Format(patcher)

	originObject.ReadFile(original)



	// Validate jmespath expression and move into path if exists
	if len(inputDiffPath) > 0 {
		_, err :=  jmespath.Compile(inputDiffPath)
		if err != nil {
			nErr := fmt.Errorf("failed to compile provided path: %T", err)
			return nErr
		}
		originObject.Data,err = jmespath.Search(inputDiffPath, originObject.Data)
		if originObject.Data == nil {
			err := errors.New("difference path returned nil object")
			return err
		} else if err != nil {
			nErr := fmt.Errorf("error pathing to object in original: %T", err)
			return nErr
		}
	}

	if len(regSkipKeys) > 0 {
		regSkip,err = regexp.Compile(regSkipKeys)
		if err != nil {
			return err
		}
	}

	newObject, err := operator.Patch(&patcher, &originObject, regSkip)
	if err != nil {
		return err
	}

	switch output {

	case "formatted":
		//writer.Write(format(objectDiff))

	case "raw":

		output, err := json.Marshal(newObject)
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


























