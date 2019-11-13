package queries

import (
    "postCorr/common"
    
    "log"
    "context"
    "strings"
    "encoding/json"
    
    "github.com/elastic/go-elasticsearch/v7" 
    "github.com/elastic/go-elasticsearch/v7/esapi"  
    
)

var es, _ = elasticsearch.NewDefaultClient()


func makeESIndexRequest(req esapi.IndexRequest) bool {
    res, err := req.Do(context.Background(),es)
    
    if err != nil {
        log.Fatalf("Error getting response: %s", err)
    }
    
    if res.IsError() {
        log.Printf("[%s] Error indexing document ID=%s", res.Status(), req.DocumentID)
        return false
    } else {
        // Deserialize the response into a map.
        var r map[string]interface{}
        if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
            log.Printf("Error parsing the response body: %s", err)
      } else {
            // Print the response status and indexed document version.
            log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
      }
        return true 
    }    
}
// Index the documents, which are split into a list of component 
func IndexDocument(indexName string, doc common.Document ) bool {
    
    var body strings.Builder
    body.WriteString(`{"components" : `)
    j, _ := json.Marshal(doc.TextComponents)
    body.WriteString(string(j))
    body.WriteString(`}`)
    
    req := esapi.IndexRequest{
        Index: indexName,
        DocumentID: doc.ID,
        Body: strings.NewReader(body.String()), 
        Refresh: "true",    
    }
    
    return makeESIndexRequest(req)
    
    
}

/**
* Puts a map with counts of occuring fingerprints into the elasticsearch index
* Mappings are converted to json to be es friendly
* Mappings have a corresponding documentID, and represent a whole document
**/

func IndexFingerPrints(indexName string, documentID string, fpCounts map[uint64]int) bool { 
    // build request body
    var body  strings.Builder
    body.WriteString(`{"fingerprints" : `)
    j, _ := json.Marshal(fpCounts)
    body.WriteString(string(j))
    body.WriteString(`}`)
    
    req := esapi.IndexRequest{
        Index: indexName,
        DocumentID: documentID,
        Body: strings.NewReader(body.String()), 
        Refresh: "true",    
    }
    
    return makeESIndexRequest(req)
 
}

 /**
func IndexAlignments(indexName string, documentID string, alignments) bool {
    
}
**/
