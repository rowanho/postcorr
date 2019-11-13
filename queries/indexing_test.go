package queries

import (
	"postCorr/fingerprinting"
    "strconv"
	"testing"
)

var testStrings = []string{
	"abcdefghijkl",
	"this is a test",
	"test unicode 切分钟",
}

/**
func TestDocumentIndexing(t *testing.T) {
    indexName := "test_fingerprints"
    
}
**/
func TestFpIndexing(t *testing.T) {
	indexName := "test_fingerprints"
	fpMaps := make([]map[uint64]int, 0)

	for _, s := range testStrings {
		fpMaps = append(fpMaps, fingerprinting.ModP(s, 2, 1))
	}

	for i, fps := range fpMaps {
        
        s := strconv.Itoa(i + 5)
		b := IndexFingerPrints(indexName, s, fps)
		if b == false {
			t.Errorf("Got error")
		}
	}
}
