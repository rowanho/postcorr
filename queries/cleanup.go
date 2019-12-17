package queries

import (
    "fmt"
)

func DeleteIndexes(indexNames []string) {
    for _, index := range indexNames {
        deleteIndex, err := es.DeleteIndex(index).Do(ctx)
        if err != nil {
            // Handle error
            fmt.Errorf("Error deleting index %s", err)
            return
        }
        if !deleteIndex.Acknowledged {
            // Not acknowledged
            fmt.Errorf("Deletion of index not acknowledged")
        } else {
            fmt.Println("Acknowledged index deletion")
        }
    }
}