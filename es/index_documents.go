package elasticlink

import (
    "log"
    "context"
    "strings"
    "strconv"
    "encoding/json"
    
    "github.com/elastic/go-elasticsearch/v7" 
    "github.com/elastic/go-elasticsearch/v7/esapi"  
)

var es, _ = elasticsearch.NewDefaultClient()

// Index the documents, which are split into a list of component strings
func IndexDocuments(indexName string, documents ) bool {
    
}

/**
* Puts a map with counts of occuring fingerprints into the elasticsearch index
* Mappings are converted to json to be es friendly
* Mappings have a corresponding documentID, and represent a whole document
**/

func IndexFingerPrints(indexName string, documentID int, fpCounts map[uint64]int) bool {
    // build request body
    var body  strings.Builder
    body.WriteString(`{"fingerprints" : `)
    j, _ := json.Marshal(fpCounts)
    body.WriteString(string(j))
    body.WriteString(`}`)
    
    req := esapi.IndexRequest{
        Index: indexName,
        DocumentID: strconv.Itoa(documentID),
        Body: strings.NewReader(body.String()), 
        Refresh: "true",    
    }
    
    res, err := req.Do(context.Background(),es)
    
    if err != nil {
        log.Fatalf("Error getting response: %s", err)
    }
    
    if res.IsError() {
        log.Printf("[%s] Error indexing document ID=%d", res.Status(), documentID)
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


func IndexAlignments(indexName string, documentID int, alignments) bool {
    
}

