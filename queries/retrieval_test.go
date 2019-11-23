package queries

import (
	"testing"	
	"time"
)



func TestDocumentRetrieval(t *testing.T) {
	for _, doc := range testDocs {
		_, err := GetDocByID(docIndexName, doc.ID)
		if err != nil {
			t.Errorf("Got error searching for document")
		}
	}
}


func TestAligmentRetrieval(t *testing.T) {
	
	// Test retrieving alignments by primary document ID
	
	primDocIds := []string{"doc1", "doc9"}
	secDocIds := []string{"doc2", "doc4"}
	
	for i, pd := range primDocIds {
		_, err := GetAlignmentsByPrimID(alignmentIndexName, pd)
		if err != nil {
			t.Errorf("Query threw an error")
		}
		
		_, err = GetAlignmentsBetween(alignmentIndexName, pd, secDocIds[i])
		if err != nil {
			t.Errorf("Query threw an error")
		} 		
	}
}


func TestLSHFingerprintRetrieval(t *testing.T) {
	time.Sleep(2 * time.Second)
	for _, doc := range testDocs {
		_, err := GetSimilarFpsLSH(fpLSHIndexName, doc.ID)
		if err != nil {
			t.Errorf("Got error searching for LSH")
		}
	}
}