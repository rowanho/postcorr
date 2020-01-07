package readWrite

import (
	"postCorr/common"
	"postCorr/queries"
	"postCorr/fingerprinting"
	"postCorr/flags"
	
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

func readAndIndex(filepath string)  error {
	
	var doc common.Document
	var err error
	if flags.FormatType == common.Plaintext {
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
	shingleLength := 2
	sampleRate := 1
	if ( flags.FpType == common.ModFP ) {
		fp := fingerprinting.ModP(string(doc.Text), shingleLength, sampleRate)
		querySuccess = queries.IndexFingerPrints(common.FpIndex, doc.ID, fp)
		if querySuccess == false {
			err =  errors.New("Couldn't index fingerprints")
		}	
	} else if (flags.FpType == common.MinhashFP) { 
		fp := fingerprinting.MinHash(filepath, string(doc.Text), shingleLength)
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
func TraverseAndIndexDocs() ([]string, error) {
	docIDs := make([]string, 0)
	if (flags.FpType == common.MinhashFP) { 
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
		fingerprinting.GetLSHObject(100, flags.JaccardThreshold, count)
	}
	err := filepath.Walk(flags.DirName,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false {
				readErr := readAndIndex(path)
				if readErr != nil {
					return readErr
				}
				docIDs = append(docIDs, path)
			}
			return nil
		},
	)
	if (flags.FpType == common.MinhashFP) { 
		fingerprinting.IndexMinHashObject()
	}
	return  docIDs, err
}