package reader

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
    
    compName := "comp0"
    newDoc := common.Document{
        ID: filepath,
        TextComponents: map[string][]rune{
            compName: ConvToStandardUnicode(text),
        },
        ComponentOrder: []string{
            compName,
        },
    }
    return newDoc, nil
}

