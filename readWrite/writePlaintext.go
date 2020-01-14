package readWrite

import (
	"fmt"
	"os"
	"strings"
)

func PlaintextWrite(docId string, text []rune) error {
	split := strings.Split(docId, "/")
	fn := split[len(split)-1]
	dirName := "corrected/" + docId[:len(docId)-len(fn)]
	os.MkdirAll(dirName, os.ModePerm)
	f, err := os.Create(dirName + fn)

	if err != nil {
		fmt.Errorf("Error, couldn't create file: %s", err)
		return err
	}

	_, err = f.WriteString(string(text))
	if err != nil {
		fmt.Errorf("Error, couldn't write to file: %s", err)
		return err
	}

	return nil
}
