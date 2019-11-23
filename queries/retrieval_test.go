package queries

import (
	"testing"	
    "fmt"
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
		alignments, err := GetAlignmentsByPrimID(alignmentIndexName, pd)
		if err != nil {
			t.Errorf("Query threw an error")
		} else if len(alignments) != 1 {
			fmt.Println(len(alignments))
			fmt.Println(alignments)
			t.Errorf("Not right number of alignments")
		}
		
		alignments, err = GetAlignmentsBetween(alignmentIndexName, pd, secDocIds[i])
		if err != nil {
			t.Errorf("Query threw an error")
		} else if len(alignments) != 1 {
			t.Errorf("Not right number of alignments")
		}		
	}
}
