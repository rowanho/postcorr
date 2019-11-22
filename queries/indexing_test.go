package queries

import (
	"postCorr/common"
	"postCorr/fingerprinting"

	"testing"
    
	"github.com/google/uuid"
)

var testDocs = []common.Document{
	{
		ID: uuid.New().String(),
		TextComponents: map[string][]rune{
			"comp1": []rune("gatctctctctctct"),
			"comp2": []rune("tctctatatatatgggg"),
			"comp3": []rune("gagagagagatctctctc"),
		},
		ComponentOrder: []string{"comp1", "comp2", "comp3"},
	},
	{
		ID: uuid.New().String(),
		TextComponents: map[string][]rune{
			"comp1": []rune("test test test test"),
			"comp2": []rune("this this this this"),
			"comp3": []rune("test this test this"),
		},
		ComponentOrder: []string{"comp1", "comp2", "comp3"},
	},
	{
		ID: uuid.New().String(),
		TextComponents: map[string][]rune{
			"comp1": []rune("test unicode"),
			"comp2": []rune("你好 再见 你好 再见"),
			"comp3": []rune("セプテンバー"),
		},
		ComponentOrder: []string{"comp1", "comp2", "comp3"},
	},
}

var testAlignments = []common.Alignment{
	{
		Score:                 5.0,
		PrimaryAl:             [][]int{{1, 2, 3, 4}, {1, 2, 3}},
		PrimaryDocumentID:     "doc1",
		PrimaryComponentIDs:   []string{"comp1", "comp2"},
		PrimaryStartComponent: "comp1",
		PrimaryEndComponent:   "comp2",
		PrimaryStartIndex:     1,
		PrimaryEndIndex:       3,
		SecondaryAl:           [][]int{{5, 6, 7, 8}, {7, 8, 10}},
		SecondaryDocumentID:   "doc2",
		SecondaryComponentIDs: []string{"comp7", "comp8"},
	},

	{
		Score:                 5.0,
		PrimaryAl:             [][]int{{7, 8}, {9, 11}, {1, 2, 3}},
		PrimaryDocumentID:     "doc9",
		PrimaryComponentIDs:   []string{"comp5", "comp6", "comp7"},
		PrimaryStartComponent: "comp5",
		PrimaryEndComponent:   "comp7",
		PrimaryStartIndex:     7,
		PrimaryEndIndex:       3,
		SecondaryAl:           [][]int{{5, 6}, {9, 13}, {7, 8, 10}},
		SecondaryDocumentID:   "doc4",
		SecondaryComponentIDs: []string{"comp7", "comp8", "comp9"},
	},
}

var docIndexName = "test_documents"
var fpIndexName = "test_fingerprints"
var alignmentIndexName = "test_alignments"

func TestDocumentIndexing(t *testing.T) {
	for _, doc := range testDocs {
		b := IndexDocument(docIndexName, doc)
		if b == false {
			t.Errorf("Got error")
		}
	}

}

func TestDocumentRetrieval(t *testing.T) {
	for _, doc := range testDocs {
		_, err := GetDocByID(docIndexName, doc.ID)
		if err != nil {
			t.Errorf("Got error searching for document")
		}
	}
}


func TestFpIndexing(t *testing.T) {

	for _, doc := range testDocs {
		s := doc.AllStrings()
		fpCounts := fingerprinting.ModP(string(s), 2, 1)
        fps := common.Fingerprints{DocumentID: doc.ID, FpCounts: fpCounts}
		b := IndexFingerPrints(fpIndexName, fps)
		if b == false {
			t.Errorf("Got error")
		}
	}

}
func TestAlignmentIndexing(t *testing.T) {

	for _, alignment := range testAlignments {
		alignmentID := uuid.New().String()
		b := IndexAlignments(alignmentIndexName, alignmentID, alignment)
		if b == false {
			t.Errorf("Got error")
		}
	}
}
