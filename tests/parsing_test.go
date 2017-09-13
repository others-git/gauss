package tests

import (



	"fmt"

	"io/ioutil"

	"testing"

	"os"
)

func ExampleParse() {

	//var JsonInput interface{}

	read, err := ioutil.ReadFile("./encoding_test.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(read)
	os.Stderr.Write(read)
}

func TestMain(*testing.M) {
	ExampleParse()
}

