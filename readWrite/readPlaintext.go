package readWrite

import (
    "postCorr/common"
    
    "io/ioutil"
    
    //"github.com/google/uuid"

)


func plaintextRead(filepath string) (common.Document, error) {
    text, err := ioutil.ReadFile(filepath)
    if err != nil {
        return common.Document{}, err
    }
    
    newDoc := common.Document{
        ID: filepath,
        Text: ConvToStandardUnicode(text),
        ComponentLengths: []int{},
    }
    return newDoc, nil
}

