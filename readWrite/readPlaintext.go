package readWrite

import (
	"postCorr/common"

	"io/ioutil"
	//"github.com/google/uuid"
)

func plaintextRead(filepath string, subpath string) (common.Document, error) {
	text, err := ioutil.ReadFile(filepath)
	if err != nil {
		return common.Document{}, err
	}

	newDoc := common.Document{
		ID:               subpath,
		Text:             ConvToStandardUnicode(text),
		ComponentLengths: []int{},
	}
	return newDoc, nil
}



func ReadString(filepath string) (string, error){
	text, err := ioutil.ReadFile(filepath)
	return string(text), err
}