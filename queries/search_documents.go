package queries

import (
    "postCorr/common"
    
    "encoding/json"
    "errors"
    "log"
    "reflect"
    "fmt"
    "strconv"
    
    "github.com/olivere/elastic/v7"     
) 

/**
* Gets a document by its document id
* This is also its id in elasticsearch
**/
func GetDocByID(indexName string, docID string) (common.Document, error) {
    
    get, err := es.Get().
        Index(indexName).
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
* Gets the similar fingerprints where we've used elastic's inbuild minhash locality sensitive hashing
* Should return a list of document IDs of the fingerprints in the same bucket 
**/

func GetSimilarFpsLSH(indexName string, documentID string) ([]string, error)  {
    get, err := es.Get().
        Index(indexName).
        Id(documentID).
        Do(ctx)
        
    if err != nil {
        fmt.Println("Couldn't find doc string")
        return []string{}, err
    }
    
    if !get.Found{
        return []string{}, errors.New("Document not found")
    } 
    
    var doc common.DocString
    json.Unmarshal(get.Source, &doc)
    fmt.Println(doc.Text)
    query := elastic.NewMoreLikeThisQuery()
    query.Field("text")
    query.MinTermFreq(1)
    query.MinDocFreq(1)
    
    it := elastic.NewMoreLikeThisQueryItem()
    it.Index(indexName)
    it.Id(doc.ID)
    
    query.LikeItems(it)
    
    res, err := es.Search().
        Index(indexName).
        Query(query).
        Pretty(true).
        Do(ctx)
        
    if err != nil {
        fmt.Println("Couldn't search for doc strings")
        return []string{}, err
    }
    if res.Hits.TotalHits.Value > 0 {
        fmt.Println("loop")
        idList := make([]string, 0)
        for _, hit := range res.Hits.Hits {
            idList = append(idList, hit.Id)
            var t common.DocString
            json.Unmarshal(hit.Source, &t)
            fmt.Println(t)
        }
        return idList, nil
    } 
    
    // No hits
    fmt.Println("no hits")
    return []string{}, nil
    
}      

/**
* Gets the alignments with similar looking priamry alignment texts 
**/
func getMatchingAlignments(indexName string, primaryID string, primaryStartIndex, 
                           primaryEndIndex, tolerance int) ([]common.Alignment, error) {
    query := elastic.NewBoolQuery()
    query = query.Must(elastic.NewTermQuery("primaryDocumentID", primaryID))
    
    simScriptStart := elastic.NewScript(`doc[primaryStartIndex].value <=  params.primStartMax ||
        || doc[PrimaryStartIndex].value >= params.primStartMin`)
    simScriptStart.Param("primStartMax", primaryStartIndex + tolerance)
    simScriptStart.Param("primStartMin", primaryStartIndex - tolerance)
    simScriptStartQuery := elastic.NewScriptQuery(simScriptStart)
        
    simScriptEnd:= elastic.NewScript(`doc[primaryEndIndex].value <=  params.primEndMax ||
        || doc[primaryEndIndex].value >= params.primEndMin`)
    simScriptStart.Param("primEndMax", primaryEndIndex + tolerance)
    simScriptStart.Param("primEndMin", primaryEndIndex - tolerance)
    simScriptEndQuery := elastic.NewScriptQuery(simScriptEnd)
        
    query = query.Must(simScriptStartQuery)
    query = query.Must(simScriptEndQuery)
    
    res, err := es.Search().
        Index(indexName).
        Query(query).
        Pretty(true).
        Do(ctx)
        
    if err != nil {
        retur
    }
    
    alignments := make([]common.Alignment, 0)
    
    var alType common.Alignment
    for _, item := range res.Each(reflect.TypeOf(alType)) {
        al := item.(common.Alignment)
        alignments = append(alignments, al)
    }
    return alignments
}