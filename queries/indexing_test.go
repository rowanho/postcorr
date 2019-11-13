package queries

import (
    "postCorr/common"
	"postCorr/fingerprinting"
    
	"testing"
)

var testDocs = []common.Document{
    {
        ID: "doc0",
        TextComponents: map[string][]rune{
            "comp1" : []rune("gatctctctctctct"),
            "comp2" :  []rune("tctctatatatatgggg"),
            "comp3" : []rune("gagagagagatctctctc"),
        },
    },
    {
        ID: "doc1",
        TextComponents: map[string][]rune{
            "comp1" : []rune("test test test test"),
            "comp2" : []rune("this this this this"),
            "comp3" : []rune("test this test this"),
        },
        
    },
    {
        ID: "doc2",
        TextComponents: map[string][]rune{
            "comp1" : []rune("test unicode"),
            "comp2" : []rune("你好 再见 你好 再见"),
            "comp3" : []rune("セプテンバー"),
        },
        
    },
    
    
}
var testStrings = []string{
	"abcdefghijkl",
	"this is a test",
	"test unicode 切分钟",
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
