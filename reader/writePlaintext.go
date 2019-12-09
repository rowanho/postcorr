package reader

import (
    "postCorr/common"
    
    "os"
    "fmt"
)


func plaintextWrite(filepath fp, doc common.Document) error {
    f, err := os.Create(fp)
    if err != nil {
        fmt.Errorf("Error, couldn't create file: %s", err)
        return err
    }
    for _, component := range doc.ComponentOrder {
        _, err := f.WriteString(string(doc.TextComponents[component]))
        if err != nil {
            fmt.Errorf("Error, couldn't write to file: %s", err)
            return err
        }
    }
    return nil
}