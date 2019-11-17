package queries

import (
    "postCorr/common"
    
    "context"
    "bytes"
    "encoding/json"
    "fmt"
    "errors"
    
    "github.com/elastic/go-elasticsearch/v7/esapi"  
    
) 


func makeESSearchRequest(indexName string,req bytes.Buffer) (*esapi.Response, error){
    
    
    res, err := es.Search(
      es.Search.WithContext(context.Background()),
      es.Search.WithIndex(indexName),
      es.Search.WithBody(&req),
      es.Search.WithTrackTotalHits(true),
      es.Search.WithPretty(),
    )
        
    if err != nil {
      eString := fmt.Sprintf("Error getting response: %s", err)
      return  &esapi.Response{}, errors.New(eString)
    }  else if res.IsError() {
        eString := fmt.Sprintf("[%s] Error indexing document", res.Status())
        return &esapi.Response{}, errors.New(eString)  
    } else {
        return res, nil
    }

    
}


func resToDoc(res map[string]interface{}) common.Document {
    e := common.Document{}
    return e
}

func GetDocByID(indexName string, docID string) (common.Document, error) {
    
    var (
        decodedRes map[string]interface{}
        req bytes.Buffer
    )
    query := map[string]interface{}{
      "query": map[string]interface{}{
        "match": map[string]interface{}{
          "_id": "docID",
        },
      },
    }
    
    json.NewEncoder(&req).Encode(query)
        
    res, err := makeESSearchRequest(indexName, req)
    defer res.Body.Close()
    
    if err != nil {
        return common.Document{}, err
    } else {
        if err := json.NewDecoder(res.Body).Decode(&decodedRes); err != nil{
            return common.Document{}, err
        }
        return resToDoc(decodedRes), nil
    }
    
}

/**
* Get the alignemts for
**/

func GetAlignmentByPrimID(indexName string, docID string) {
    
}

/**

* Gets the alignemnts between two documents
**/

func GetAlignmentsBetween(indexName string, primaryID string, secondaryID string) {
    
}


