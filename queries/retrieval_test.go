package queries

import (
	"fmt"
	"testing"
	"time"
)

func TestDocumentRetrieval(t *testing.T) {
	time.Sleep(1 * time.Second)
	for _, doc := range testDocs {
		_, err := GetDocByID(docIndexName, doc.ID)
		if err != nil {
			t.Errorf("Got error searching for document")
		}
	}
}

func TestAligmentRetrieval(t *testing.T) {
	time.Sleep(1 * time.Second)
	// Test retrieving alignments by primary document ID

	primDocIds := []string{"doc1", "doc9"}
	secDocIds := []string{"doc2", "doc4"}

	for i, pd := range primDocIds {
		als, err := GetAlignmentsByPrimID(alignmentIndexName, pd)
		if err != nil {
			t.Errorf("Query threw an error")
		}

		for _, al := range als {
			matching, err := GetMatchingAlignments(alignmentIndexName,
				al, 3)
			if err != nil {
				t.Errorf("Matching query threw an error %s: ", err)
			}

			fmt.Printf("Alignment with primary ID %s matched with %d others. \n", pd, len(matching))
		}

		_, err = GetAlignmentsBetween(alignmentIndexName, pd, secDocIds[i])
		if err != nil {
			t.Errorf("Query threw an error")
		}
	}
}

func TestLSHFingerprintRetrieval(t *testing.T) {
	time.Sleep(1 * time.Second)
	for _, doc := range testDocs {
		_, err := GetSimilarFpsLSH(fpLSHIndexName, doc.ID)
		if err != nil {
			t.Errorf("Got error searching for LSH: %s", err)
		}
	}
}
