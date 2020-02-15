package readWrite

import (
	"postCorr/flags"
	
	"fmt"
	"os"
	"strings"
	"strconv"
	"path"
)
var first = true

func PlaintextWrite(docId string, text []rune) error {

	if first {
		i := 0
		newDir := flags.OutDir
		_, err := os.Stat(newDir)
		for  !os.IsNotExist(err){
			i +=  1
			newDir =  flags.OutDir + strconv.Itoa(i)
			_, err = os.Stat(newDir);
		}
		first = false
		flags.OutDir = newDir
	}
	
	split := strings.Split(docId, "/")
	fn := split[len(split)-1]
	dirName := path.Join(flags.OutDir, docId[:len(docId)-len(fn)])
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
