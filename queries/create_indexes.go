package queries

import (
    "strconv"
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


/**
* Creates the document index with minhash, a LSH bucket hasher
**/
func CreateLSHFingerprintIndex(indexName string, shingleMin int, shingleMax int, noBuckets int) error {
    
    mappings := `{
      "settings": {
        "analysis": {
          "filter": {
            "shingler": { 
              "type": "shingle",
              "min_shingle_size": `+ strconv.Itoa(shingleMin) +`,
              "max_shingle_size": `+ strconv.Itoa(shingleMax) +`,
              "output_unigrams": false
            },
            "minhasher": {
              "type": "min_hash",
              "hash_count": 1,   
              "bucket_count": `+ strconv.Itoa(noBuckets) +`, 
              "hash_set_size": 1, 
              "with_rotation": true 
            }
          },
          "analyzer": {
            "fingerprinting": {
              "tokenizer": "standard",
              "filter": [
                "shingler",
                "minhasher"
              ]
            }
          }
        }
      },
      "mappings": {
        "properties": {
          "text": {
            "type": "text",
            "analyzer": "fingerprinting"
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