package queries

import (
	"postCorr/common"

	"context"
	"fmt"
	"log"
    
	"github.com/olivere/elastic/v7" 
)

var es, _ = elastic.NewClient()
var ctx = context.Background()


// Index the documents, which are split into a list of component
func IndexDocument(indexName string, doc common.Document) bool {

    put, err := es.Index().
        Index(indexName).
        Id(doc.ID).
        BodyJson(doc).
        Do(ctx)
    
    if err != nil {
        log.Printf("Error indexing document: %s", err)  
        return false      
    }
    fmt.Printf("Indexed document %s to index %s\n", put.Id, put.Index)
    return true
}

/**
* Puts a map with counts of occuring fingerprints into the elasticsearch index
* Mappings are converted to json to be es friendly
* Mappings have a corresponding documentID, and represent a whole document
**/

func IndexFingerPrints(indexName string, fps common.Fingerprints) bool {

    put, err := es.Index().
        Index(indexName).
        Id(fps.DocumentID).
        BodyJson(fps).
        Do(ctx)
    
    if err != nil {
        log.Printf("Error indexing alignment: %s", err)  
        return false      
    }
    fmt.Printf("Indexed alignment %s to index %s\n", put.Id, put.Index)
    return true
}

func IndexFingerPrintsForLSH(indexName string, fps common.DocString) bool {
	put, err := es.Index().
        Index(indexName).
        Id(fps.ID).
        BodyJson(fps).
        Do(ctx)
    
    if err != nil {
        log.Printf("Error indexing alignment: %s", err)  
        return false      
    }
    fmt.Printf("Indexed alignment %s to index %s\n", put.Id, put.Index)
    return true
	
}

/**
* Puts an alignment into elasticsearch
* The id field is made up of the priamry document ID, component ID, and start/end indices
*
*
**/

func IndexAlignments(indexName string, alignment common.Alignment) bool {

    put, err := es.Index().
        Index(indexName).
        Id(alignment.ID).
        BodyJson(alignment).
        Do(ctx)
    
    if err != nil {
        log.Printf("Error indexing alignment: %s", err)  
        return false      
    }
    fmt.Printf("Indexed alignment %s to index %s\n", put.Id, put.Index)
    return true
}
