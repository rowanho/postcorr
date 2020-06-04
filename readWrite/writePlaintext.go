package readWrite

import (
	"postCorr/common"

	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	first = true
	outdir = common.OutDir
)

func PlaintextWrite(docId string, text []rune) error {

	if first {
		i := 0
		newDir := common.OutDir
		_, err := os.Stat(newDir)
		for !os.IsNotExist(err) {
			i += 1
			newDir = common.OutDir + strconv.Itoa(i)
			_, err = os.Stat(newDir)
		}
		first = false
		outdir = newDir
	}

	split := strings.Split(docId, "/")
	fn := split[len(split)-1]
	dirName := path.Join(outdir, docId[:len(docId)-len(fn)])
	os.MkdirAll(dirName, os.ModePerm)
	f, err := os.Create(path.Join(dirName, fn))

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
