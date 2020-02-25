package readWrite

import (
	"postCorr/common"
	"postCorr/flags"

	"os"
	"path/filepath"

	"golang.org/x/text/unicode/norm"
)

/**
* Converts the bytes to a rune array, based off of our standardised unicode encoding
* Here we ae using the NFC canonical equivalence
**/

func ConvToStandardUnicode(b []byte) []rune {
	return []rune(string(norm.NFC.Bytes(b)))
}

/**
* Traverses the dataset folder and indexes the document
* Returns a list of the document names, and an error
**/
func TraverseDocs() ([]common.Document, error) {
	docs := make([]common.Document, 0)
	err := filepath.Walk(flags.DirName,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false {
				subpath := path[len(flags.DirName) + 1:]
				doc, readErr := plaintextRead(path, subpath)
				docs = append(docs, doc)
				
				if readErr != nil {
					return readErr
				}
			}
			return nil
		},
	)
	return docs, err
}
