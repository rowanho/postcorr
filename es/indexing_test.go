package elasticlink

import (
    "testing"
    "postCorr/fingerprinting"
)

var testStrings = []string {
    "abcdefghijkl",
    "this is a test",
    "test unicode 切分钟",
}
func TestFpIndexing(t *testing.T) {
    indexName := "test_fingerprints"
    fpMaps := make([]map[uint64]int, 0)
    
    for _, s := range testStrings {
        fpMaps = append(fpMaps, fingerprinting.ModP(s, 2, 1))
    }
    
    for i, fps := range fpMaps {
        
        b := IndexFingerPrints(indexName, i, fps)
		if b == false {
			t.Errorf("Got error")
		}
	}
}