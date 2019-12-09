package readWrite

import (
    "postCorr/common"
    
    "os"
    "fmt"
)


func PlaintextWrite(fp string, doc common.Document) error {
    f, err := os.Create(fp)
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