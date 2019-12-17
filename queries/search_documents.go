package queries

import (
    "postCorr/common"
    "postCorr/fingerprinting"
    
    "encoding/json"
    "errors"
    "log"
    "reflect"
    
    "github.com/olivere/elastic/v7"    
    "github.com/DearMadMan/minhash"
 
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
        return common.Alignment{}, err
    }
    
    if get.Found{
        var al common.Alignment
        json.Unmarshal(get.Source, &al)
        return al, nil
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


// If the jaccard score is over the threshold, add it to the map as a similar document
func GetSimilarFps(indexName string, targetDocumentID string, docIDList [] string, jaccardThreshold float64) (map[string]bool, error) {
    
    similarDocIDs := map[string]bool{}
    var targetFp map[uint64]int
    
    get, err := es.Get().
    Index(indexName).
    Id(targetDocumentID).
    Do(ctx)
    if err != nil {
        return similarDocIDs, err;
    } else{
        json.Unmarshal(get.Source, &targetFp)    
    }

    for _, docID := range docIDList {
        if docID == targetDocumentID {
            continue;
        }
        get, err := es.Get().
        Index(indexName).
        Id(docID).
        Do(ctx)
        if err != nil {
            continue;
        }
        var fp map[uint64]int
        json.Unmarshal(get.Source, &fp)
        if fingerprinting.FpJaccardScore(targetFp, fp) > jaccardThreshold {
            similarDocIDs[docID] = true
        }
    }
    return similarDocIDs, nil
}


func GetSimilarMinHashes(indexName string, targetDocumentID string, docIDList []string, jaccardThreshold float64) (map[string]bool, error) {
    similarDocIDs := map[string]bool{}
    var targetFp minhash.Set
    
    get, err := es.Get().
    Index(indexName).
    Id(targetDocumentID).
    Do(ctx)
    if err != nil {
        return similarDocIDs, err;
    } else{
        json.Unmarshal(get.Source, &targetFp)    
    }

    for _, docID := range docIDList {
        if docID == targetDocumentID {
            continue;
        }
        get, err := es.Get().
        Index(indexName).
        Id(docID).
        Do(ctx)
        if err != nil {
            continue;
        }
        var fp minhash.Set
        json.Unmarshal(get.Source, &fp)
        if targetFp.Jaccard(fp) > jaccardThreshold {
            similarDocIDs[docID] = true
        }
    }
    return similarDocIDs, nil    
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