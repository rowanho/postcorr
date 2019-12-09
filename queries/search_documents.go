package queries

import (
    "postCorr/common"
    
    "encoding/json"
    "errors"
    "log"
    "reflect"
    
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
* Gets a document by its document id
* This is also its id in elasticsearch
**/
func GetAlignmentByID(indexName string, alID string) (common.Alignment, error) {
    
    get, err := es.Get().
        Index(indexName).
        Id(alID).
        Do(ctx)
    if err != nil {
        log.Printf("Error getting document: %s", err)
        return common.Document{}, err
    }
    
    if get.Found{
        var doc common.Alignment
        json.Unmarshal(get.Source, &doc)
        return doc, nil
    } else {
        return common.Alignment{}, errors.New("Alignment not found")
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
        log.Println("Couldn't find doc string")
        return []string{}, err
    }
    
    if !get.Found{
        return []string{}, errors.New("Document not found")
    } 
    
    var doc common.DocString
    json.Unmarshal(get.Source, &doc)
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
        return []string{}, err
    }
    if res.Hits.TotalHits.Value > 0 {
        idList := make([]string, 0)
        for _, hit := range res.Hits.Hits {
            idList = append(idList, hit.Id)
        }
        return idList, nil
    } 
    
    // No hits
    return []string{}, nil
    
}      

/**
* Gets the alignments where the primary alignment region is similar to the primary
* alignment region of our query alignment
**/
func GetMatchingAlignments(indexName string, al common.Alignment, tolerance int) ([]string, error) {
    query := elastic.NewBoolQuery()
    query = query.Must(elastic.NewTermQuery("primaryDocumentID", al.PrimaryDocumentID))
    query = query.MustNot(elastic.NewTermQuery("_id", al.ID))
    
    simScriptStart := elastic.NewScript(`doc['primaryStartIndex'].value <=  params.startMax
        || doc['primaryStartIndex'].value >= params.startMin`)
    simScriptStart.Param("startMax", al.PrimaryStartIndex + tolerance)
    simScriptStart.Param("startMin", al.PrimaryStartIndex - tolerance)
    simScriptStartQuery := elastic.NewScriptQuery(simScriptStart)
        
    simScriptEnd:= elastic.NewScript(`doc['primaryEndIndex'].value <=  params.endMax
        || doc['primaryEndIndex'].value >= params.endMin`)
    simScriptEnd.Param("endMax", al.PrimaryEndIndex + tolerance)
    simScriptEnd.Param("endMin", al.PrimaryEndIndex - tolerance)
    simScriptEndQuery := elastic.NewScriptQuery(simScriptEnd)
        
    query = query.Filter(simScriptStartQuery)
    query = query.Filter(simScriptEndQuery)
    
    res, err := es.Search().
        Index(indexName).
        Query(query).
        Pretty(true).
        Do(ctx)
        
    if err != nil {
        log.Printf("Error getting similar alignments: %s", err)        
        return []string{}, err
    }
    
    if res.Hits.TotalHits.Value > 0 {
        idList := make([]string, 0)
        for _, hit := range res.Hits.Hits {
            idList = append(idList, hit.Id)
        }
        log.Println(idList)
        return idList, nil
    } 
    log.Println("empty!")
    return []string{}, nil
    
}

/**
* Gets the alignments where the primary alignment region is similar to the secondary
* alignment region of our query alignment
**/

func GetConnectedAlignments(indexName string, al common.Alignment, tolerance int) ([]string, error) {
    query := elastic.NewBoolQuery()
    query = query.Must(elastic.NewTermQuery("primaryDocumentID", al.SecondaryDocumentID))
    query = query.MustNot(elastic.NewTermQuery("_id", al.ID))
    query = query.MustNot(elastic.NewTermQuery("secondaryDocumentID", al.PrimaryDocumentID))

    simScriptStart := elastic.NewScript(`doc['primaryStartIndex'].value <=  params.startMax
        || doc['primaryStartIndex'].value >= params.startMin`)
    simScriptStart.Param("startMax", al.SecondaryStartIndex + tolerance)
    simScriptStart.Param("startMin", al.SecondaryStartIndex - tolerance)
    simScriptStartQuery := elastic.NewScriptQuery(simScriptStart)
        
    simScriptEnd:= elastic.NewScript(`doc['primaryEndIndex'].value <=  params.endMax
        || doc['primaryEndIndex'].value >= params.endMin`)
    simScriptEnd.Param("endMax", al.SecondaryEndIndex + tolerance)
    simScriptEnd.Param("endMin", al.SecondaryEndIndex - tolerance)
    simScriptEndQuery := elastic.NewScriptQuery(simScriptEnd)
        
    query = query.Filter(simScriptStartQuery)
    query = query.Filter(simScriptEndQuery)

    res, err := es.Search().
        Index(indexName).
        Query(query).
        Pretty(true).
        Do(ctx)
        
    if err != nil {
        log.Printf("Error getting similar alignments: %s", err)        
        return []string{}, err
    }

    if res.Hits.TotalHits.Value > 0 {
        idList := make([]string, 0)
        for _, hit := range res.Hits.Hits {
            idList = append(idList, hit.Id)
        }
        log.Println(idList)
        return idList, nil
    } 
    log.Println("empty!")
    return []string{}, nil
}