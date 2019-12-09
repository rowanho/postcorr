package readWrite

import (
    "postCorr/common"
    
    "os"
    "fmt"
    "strings"
)


func PlaintextWrite(fp string, doc common.Document) error {
    split := strings.Split(fp, "/")
    fn := split[len(split) - 1]
    dirName  := "corrected/" + fp[:len(fp) - len(fn)]
    os.MkdirAll(dirName, os.ModePerm)
    f, err := os.Create(dirName + fn)
    
    if err != nil {
        fmt.Errorf("Error, couldn't create file: %s", err)
        return err
    }
    _, err = f.WriteString(string(doc.Text))
    if err != nil {
        fmt.Errorf("Error, couldn't write to file: %s", err)
        return err
    }
    
    return nil
}