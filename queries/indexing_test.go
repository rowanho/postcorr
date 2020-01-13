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
			"comp1": []rune("test test test test"),
			"comp2": []rune("this this thi this"),
			"comp3": []rune("test this tet this"),
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
		ID:                  uuid.New().String(),
		Score:               5.0,
		PrimaryAl:           []int{1, 2, 3, 4},
		PrimaryDocumentID:   "doc1",
		PrimaryStartIndex:   1,
		PrimaryEndIndex:     4,
		SecondaryAl:         []int{5, 6, 7, 8, 10},
		SecondaryDocumentID: "doc2",
		SecondaryStartIndex: 5,
		SecondaryEndIndex:   10,
	},

	{
		ID:                  uuid.New().String(),
		Score:               5.0,
		PrimaryAl:           []int{7, 8, 9, 11},
		PrimaryDocumentID:   "doc9",
		PrimaryStartIndex:   7,
		PrimaryEndIndex:     11,
		SecondaryAl:         []int{5, 6, 9, 13},
		SecondaryDocumentID: "doc4",
		SecondaryStartIndex: 5,
		SecondaryEndIndex:   13,
	},
}

var docIndexName = "test_documents"
var fpIndexName = "test_fingerprints"
var fpLSHIndexName = "test_lsh_fingerprints"
var alignmentIndexName = "test_alignments"

func TestDocumentIndexing(t *testing.T) {
	for _, doc := range testDocs {
		b := IndexDocument(docIndexName, doc)
		if b == false {
			t.Errorf("Got error")
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
		b := IndexAlignments(alignmentIndexName, alignment)
		if b == false {
			t.Errorf("Got error")
		}
	}
}

func TestLSHFpIndexing(t *testing.T) {

	numBuckets := 10
	shingleMin := 2
	shingleMax := 3

	err := CreateLSHFingerprintIndex(fpLSHIndexName, shingleMin, shingleMax, numBuckets)
	if err != nil {
		t.Errorf("Error creating mappings")
	}

	// Index DocStrings , which get hashed by elasticsearch
	for _, doc := range testDocs {
		docText := string(doc.AllStrings())
		docString := common.DocString{
			ID:   doc.ID,
			Text: docText,
		}
		b := IndexFingerPrintsForLSH(fpLSHIndexName, docString)
		if b == false {
			t.Errorf("Error indexing fingerprint with LSH")
		}

	}

}
