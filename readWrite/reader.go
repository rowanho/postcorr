package readWrite

import (
	"postCorr/common"
	"postCorr/fingerprinting"
	"postCorr/flags"
	"postCorr/queries"

	"errors"
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
func TraverseDocs() ([]string, error) {
	docs := make([]common.Document, 0)
	if flags.FpType == common.MinhashFP {
		count := 0
		err := filepath.Walk(flags.DirName,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				count += 1
				return nil
			},
		)
		if err != nil {
			return docIDs, err
		}
	}

	err := filepath.Walk(flags.DirName,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false {
				var readErr error
				if flags.FormatType == common.Plaintext {
					doc, readErr = plaintextRead(filepath)
				}
				if readErr != nil {
					return readErr
				}
				docs = append(docs, doc)
			}
			return nil
		},
	)
	return docs, err
}
