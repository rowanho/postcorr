package queries

import (
    "postCorr/common"
    
    "encoding/json"
    "errors"
    "log"
    "reflect"
    "fmt"
    
    "github.com/olivere/elastic/v7"     
) 

/**
* Gets a document by its document id
* This is also its id in elasticsearch
**/
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
* Retrieves alignments that have the same primary id
**/

func GetAlignmentsByPrimID(indexName string, primID string) ([]common.Alignment, error) {
    query := elastic.NewTermQuery("primaryDocumentID", primID)
    
    src, err := query.Source()
    if err != nil {
      panic(err)
    } else {
        fmt.Println(src)
    }
    
    res, err := es.Search().
        Index(indexName).
        Query(query).
        Pretty(true).
        Do(ctx)
        
    if err != nil {
        log.Printf("Error searching for alignments: %s", err)
        return []common.Alignment{}, err
    }
    
    alignments := make([]common.Alignment, 0)
    var alType common.Alignment
    for _, item := range res.Each(reflect.TypeOf(alType)) {
        al := item.(common.Alignment)
        alignments = append(alignments, al)
    }
    
    return alignments, nil
}    

 /**
 * Retrieves alignments with the same primary and secondary id
 **/

func GetAlignmentsBetween(indexName string, primaryID string, secondaryID string) ([]common.Alignment, error) {
    query := elastic.NewBoolQuery()
    query = query.Must(elastic.NewTermQuery("primaryDocumentID", primaryID))
    query = query.Must(elastic.NewTermQuery("secondaryDocumentID", secondaryID))
    
    src, err := query.Source()
    if err != nil {
      panic(err)
    } else {
        fmt.Println(src)
    }
    
    res, err := es.Search().
        Index(indexName).
        Query(query).
        Pretty(true).
        Do(ctx)
    
    if err != nil {
        log.Printf("Error searching for alignments: %s", err)
        return []common.Alignment{}, err
    }
    
    alignments := make([]common.Alignment, 0)
    var alType common.Alignment
    for _, item := range res.Each(reflect.TypeOf(alType)) {
        al := item.(common.Alignment)
        alignments = append(alignments, al)
    }
    
    return alignments, nil
}

