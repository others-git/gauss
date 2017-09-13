package parsing

import (
	/*
	"encoding/json"
	"github.com/go-yaml/yaml"

	"reflect"AT
	*/
)



type readerType struct {}

func (r *readerType) yamlMarshal() interface{} {
	var placeholder interface{}

	return placeholder
}
func (r *readerType) yamlUnmarshal() interface{} {
	var placeholder interface{}

	return placeholder
}

func (r *readerType) jsonUnmarshal() interface{} {
	var placeholder interface{}

	return placeholder
}
func (r *readerType) jsonMarshal() interface{} {
	var placeholder interface{}

	return placeholder
}


func Detector() readerType{
	var reader readerType

	return reader
}



