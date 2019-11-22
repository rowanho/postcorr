package queries

import (
    "postCorr/common"
    
    "encoding/json"
    "errors"
    "log"
        
) 

func GetDocByID(indexName string, docID string) (common.Document, error) {
    
    get, err := es.Get().
        Index(indexName).
        Type("document").
        Id(docID).
        Do(ctx)
    if err != nil {
        log.Printf("Error getting document: %s", err)
        return common.Document{}, err
    }
    
    if get.Found{
        var doc common.Document
        json.Unmarshal(get.Source, &doc)
        return doc, nil
    } else {
        return common.Document{}, errors.New("Document not found")
    }
    
    
}

/**

func GetAlignmentByPrimID(indexName string, docID string) (common.Alignment, error) {



* Gets the alignemnts between two documents

func GetAlignmentsBetween(indexName string, primaryID string, secondaryID string) {
    
}
**/

