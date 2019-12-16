package queries

import (
    "errors"
    "fmt"
)

/**
* Creates the alignment index
**/
func CreateAlignmentIndex(indexName string) error {
    mappings := `{
        "mappings": {
            "properties": {
                "primaryDocumentID": {
                    "type": "keyword"
                },
                "secondaryDocumentID": {
                    "type": "keyword"
                }
            }
        }
            
    }`
    createIndex, err := es.CreateIndex(indexName).  
        BodyString(mappings).
        Do(ctx)
        
    if err != nil {
        fmt.Printf("Error creating mappings: %s", err)
        return err
    }
    
    if !createIndex.Acknowledged {
        fmt.Println("Error: Index creation not acknowledged")
        return errors.New("Index creation not acknowledged")
    } else {
        return nil
    }
    
}

func CreateFingerprintIndex(indexName string) error {
    mappings := `{
        "settings": {
            "index.mapping.total_fields.limit": 100000 
        }
    }`
    createIndex, err := es.CreateIndex(indexName).  
        BodyString(mappings).
        Do(ctx)
        
    if err != nil {
        fmt.Printf("Error creating mappings: %s", err)
        return err
    }
    
    if !createIndex.Acknowledged {
        fmt.Println("Error: Index creation not acknowledged")
        return errors.New("Index creation not acknowledged")
    } else {
        return nil
    }
    
}


