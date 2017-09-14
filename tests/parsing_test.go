package tests

import (



	"fmt"


	"testing"

	"github.com/dimchansky/utfbom"
	"io/ioutil"
	"bytes"
	"os"
)

func ExampleParse() {

	//var JsonInput interface{}



	f, err := ioutil.ReadFile("./origin.json")
	if err != nil {
		fmt.Println(err)
	}

	o,_ := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(f)))
	os.Stdout.Write(o)


}

func TestMain(*testing.M) {
	ExampleParse()
}

