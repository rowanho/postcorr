package readWrite

import (
	"postCorr/common"
	"postCorr/queries"
	"postCorr/fingerprinting"
	
	"path/filepath"
	"os"
	"errors"
	
	"golang.org/x/text/unicode/norm"
)

/**
* Converts the bytes to a rune array, based off of our standardised unicode encoding
* Here we ae using the NFC canonical equivalence
**/

func ConvToStandardUnicode(b []byte) []rune {
	return []rune(string(norm.NFC.Bytes(b)))
}

func readAndIndex(filepath string, formatType string, fpType string)  error {
	
	var doc common.Document
	var err error
	if formatType == common.Plaintext {
		doc, err = plaintextRead(filepath)
	}
	
	if err != nil {
		return err
	}
	
	querySuccess := queries.IndexDocument(common.DocumentIndex, doc)
	
	if querySuccess == false {
		return errors.New("Couldn't index document")
	}
	
	err = nil
	if ( fpType == common.ModFP ) {
		fp := fingerprinting.ModP(string(doc.Text), 7 , 10)
		querySuccess = queries.IndexFingerPrints(common.FpIndex, doc.ID, fp)
		if querySuccess == false {
			err =  errors.New("Couldn't index fingerprints")
		}	
	} else if (fpType == common.MinhashFP) { 
		fp := fingerprinting.MinHash(string(doc.Text))
		querySuccess = queries.IndexMinhash(common.MinHashIndex, doc.ID, fp)
		if querySuccess == false {
			err =  errors.New("Couldn't index fingerprints")
		}	
	}
	
	return err
}



/**
* Traverses the dataset folder and indexes the document
* Returns a list of the document names, and an error
**/
func TraverseAndIndexDocs(dirName string, formatType string, fpType string) ([]string, error) {
	docIDs := make([]string, 0)
	err := filepath.Walk(dirName,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false {
				readErr := readAndIndex(path, formatType, fpType)
				if readErr != nil {
					return readErr
				}
				docIDs = append(docIDs, path)
			}
			return nil
		})
	
	return  docIDs, err
}