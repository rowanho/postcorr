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
            "comp1" : []rune("gatctctctctctct"),
            "comp2" :  []rune("tctctatatatatgggg"),
            "comp3" : []rune("gagagagagatctctctc"),
        },
        ComponentOrder: []string{"comp1", "comp2", "comp3"},
    },
    {
        ID: uuid.New().String(),
        TextComponents: map[string][]rune{
            "comp1" : []rune("test test test test"),
            "comp2" : []rune("this this this this"),
            "comp3" : []rune("test this test this"),
        },
        ComponentOrder: []string{"comp1", "comp2", "comp3"},
        
    },
    {
        ID: uuid.New().String(),
        TextComponents: map[string][]rune{
            "comp1" : []rune("test unicode"),
            "comp2" : []rune("你好 再见 你好 再见"),
            "comp3" : []rune("セプテンバー"),
        },
        ComponentOrder: []string{"comp1", "comp2", "comp3"},
        
    },
    
    
}

var testAlignments = []common.AlignmentCluster{
    {
         PrimaryAl : [][]int{{1,2,3,4},  {1,2,3},},
         PrimaryDocumentID : "doc1",
         PrimaryComponentIDs : []string{"comp1", "comp2",},
         PrimaryStartComponent : "comp1",
         PrimaryEndComponent : "comp2",
         PrimaryStartIndex : 1,
         PrimaryEndIndex : 3,
         SecondaryAl : [][]int{{5,6,7,8}, {7,8,10},},
         SecondaryDocumentID : "doc2",
         SecondaryComponentIDs : []string{"comp7", "comp8",},
    },
    
    {
        PrimaryAl : [][]int{{7,8}, {9,11}, {1,2,3}},
        PrimaryDocumentID : "doc9",
        PrimaryComponentIDs : []string{"comp5", "comp6", "comp7"},
        PrimaryStartComponent : "comp5",
        PrimaryEndComponent : "comp7",
        PrimaryStartIndex : 7,
        PrimaryEndIndex : 3,
        SecondaryAl : [][]int{{5,6}, {9,13}, {7,8,10}},
        SecondaryDocumentID :  "doc4",
        SecondaryComponentIDs : []string{"comp7", "comp8", "comp9"},
        
    },
    
}


func TestDocumentIndexing(t *testing.T) {
    indexName := "test_documents"
    for _, doc := range testDocs {
        IndexDocument(indexName, doc)
    }
    
}

func TestFpIndexing(t *testing.T) {
	indexName := "test_fingerprints"

	for _, doc := range testDocs {
        s := doc.AllStrings()
		fps :=  fingerprinting.ModP(string(s), 2, 1)
		b := IndexFingerPrints(indexName, doc.ID, fps)
        if b == false {
			t.Errorf("Got error")
		}
	}

}

func TestAlignmentIndexing(t *testing.T) {
    indexName := "test_alignments"
    
    for _, alignment := range testAlignments {
        docID := uuid.New().String()
        b := IndexAlignments(indexName, docID, alignment)
        if b == false {
            t.Errorf("Got error")
        }
    }
}
