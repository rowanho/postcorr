package readWrite

import (
    "postCorr/common"
    "postCorr/queries"
    
    "os"
    "fmt"
    "strings"
)


func PlaintextWrite(docId string) error {
    split := strings.Split(docId, "/")
    fn := split[len(split) - 1]
    dirName  := "corrected/" + docId[:len(docId) - len(fn)]
    os.MkdirAll(dirName, os.ModePerm)
    f, err := os.Create(dirName + fn)
    
    if err != nil {
        fmt.Errorf("Error, couldn't create file: %s", err)
        return err
    }
    
    doc, _ := queries.GetDocByID(common.DocumentIndex, docId)
    _, err = f.WriteString(string(doc.Text))
    if err != nil {
        fmt.Errorf("Error, couldn't write to file: %s", err)
        return err
    }
    
    return nil
}