package queries

import (
	"postCorr/common"

	"context"
	"fmt"
	"log"
  "time"
	
	"github.com/olivere/elastic/v7" 
)

var es, _ = elastic.NewClient(elastic.SetSnifferTimeout(10 * time.Second),elastic.SetHealthcheckInterval(10 * time.Second))
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
*  Bulk indexes documents
*  Includes a refresh parameter, which waits for indexing
**/

func BulkUpdateDocuments(indexName string, docs map[string][]rune) bool {
		bulkRequest := es.Bulk()
		for id, text := range docs {
				req := elastic.NewBulkUpdateRequest().
				Index(indexName).
				Id(id).                
				Doc(struct {
				  Text []rune `json:"text"`
				}{
				  Text: text,
				})

				bulkRequest = bulkRequest.Add(req)
		}
		_, err := bulkRequest.Refresh("wait_for").Do(ctx)
		
		if err != nil {
			log.Printf("Error updating documents: %s", err)  
			return false
		} else {
			return true
		}
}
/**
* Puts a map with counts of occuring fingerprints into the elasticsearch index
* Mappings are converted to json to be es friendly
* Mappings have a corresponding documentID, and represent a whole document
**/

func IndexFingerPrints(indexName string, docID string, fp map[uint64]int) bool {

    put, err := es.Index().
        Index(indexName).
        Id(docID).
        BodyJson(fp).
        Do(ctx)
    
    if err != nil {
        log.Printf("Error indexing fingerprint: %s", err)  
        return false      
    }
    fmt.Printf("Indexed fingerprint %s to index %s\n", put.Id, put.Index)
    return true
}

/**
* Puts a minhash into elasticsearch
**/
func IndexMinhash(indexName string, docID string, fp common.LSH_fp) bool {
	put, err := es.Index().
        Index(indexName).
        Id(docID).
        BodyJson(fp).
        Do(ctx)
    
    if err != nil {
        log.Printf("Error indexing fingerprint: %s", err)  
        return false      
    }
    fmt.Printf("Indexed fingerprint %s to index %s\n", put.Id, put.Index)
    return true
}

/**
* Puts an alignment into elasticsearch
* The id field is made up of the priamry document ID, component ID, and start/end indices
*
*
**/

func IndexAlignment(indexName string, alignment common.Alignment) bool {

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
