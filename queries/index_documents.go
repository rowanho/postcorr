package queries

import (
	"postCorr/common"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var es, _ = elasticsearch.NewDefaultClient()

func makeESIndexRequest(req esapi.IndexRequest) bool {
	res, err := req.Do(context.Background(), es)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), req.DocumentID)
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

// Index the documents, which are split into a list of component
func IndexDocument(indexName string, doc common.Document) bool {

	var body strings.Builder
	body.WriteString(`{"components" : {`)
	last := len(doc.ComponentOrder) - 1
	for i, compID := range doc.ComponentOrder {
		body.WriteString(`"` + compID + `"` + ` : `)
		body.WriteString(`"` + string(doc.TextComponents[compID]) + `"`)
		if i < last {
			body.WriteString(` ,`)
		}
	}
	body.WriteString(`},`)
	body.WriteString(`"componentOrder" : `)
	j, _ := json.Marshal(doc.ComponentOrder)
	body.WriteString(string(j))
	body.WriteString(`}`)

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: doc.ID,
		Body:       strings.NewReader(body.String()),
		Refresh:    "true",
	}

	return makeESIndexRequest(req)

}

/**
* Puts a map with counts of occuring fingerprints into the elasticsearch index
* Mappings are converted to json to be es friendly
* Mappings have a corresponding documentID, and represent a whole document
**/

func IndexFingerPrints(indexName string, documentID string, fpCounts map[uint64]int) bool {
	// build request body
	var body strings.Builder
	body.WriteString(`{"fingerprints" : `)
	j, _ := json.Marshal(fpCounts)
	body.WriteString(string(j))
	body.WriteString(`}`)

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       strings.NewReader(body.String()),
		Refresh:    "true",
	}

	return makeESIndexRequest(req)

}

/**
* Puts an alignment into elasticsearch
* The id field is made up of the priamry document ID, component ID, and start/end indices
*
*
**/

func IndexAlignments(indexName string, docID string, alignment common.Alignment) bool {
	var body strings.Builder
	body.WriteString(`{"primaryDocumentID" : `)
	body.WriteString(`"` + alignment.PrimaryDocumentID + `"`)
	body.WriteString(`, `)

	body.WriteString(`"primaryAlignmentIndices" : `)
	body.WriteString(`{`)
	last := len(alignment.PrimaryComponentIDs) - 1
	for i, compID := range alignment.PrimaryComponentIDs {
		body.WriteString(`"` + compID + `" : `)
		a, _ := json.Marshal(alignment.PrimaryAl[i])
		if i < last {
			body.WriteString(string(a) + `,`)
		} else {
			body.WriteString(string(a))
		}

	}
	body.WriteString(`}, `)

	body.WriteString(`"primaryStartIndex" : `)
	body.WriteString(strconv.Itoa(alignment.PrimaryStartIndex))
	body.WriteString(`, `)

	body.WriteString(`"primaryEndIndex" : `)
	body.WriteString(strconv.Itoa(alignment.PrimaryEndIndex))
	body.WriteString(`, `)

	body.WriteString(`"score" : `)
	body.WriteString(fmt.Sprintf("%f", alignment.Score))
	body.WriteString(`, `)

	body.WriteString(`"primaryStartComponent" : `)
	body.WriteString(`"` + alignment.PrimaryStartComponent + `"`)
	body.WriteString(`, `)

	body.WriteString(`"primaryEndComponent" : `)
	body.WriteString(`"` + alignment.PrimaryEndComponent + `"`)
	body.WriteString(`, `)

	body.WriteString(`"secondaryDocumentID" : `)
	body.WriteString(`"` + alignment.SecondaryDocumentID + `"`)
	body.WriteString(`, `)

	body.WriteString(`"secondaryAlignmentIndices" : `)
	body.WriteString(`{`)

	last = len(alignment.SecondaryComponentIDs) - 1
	for i, compID := range alignment.SecondaryComponentIDs {
		body.WriteString(`"` + compID + `" : `)
		a, _ := json.Marshal(alignment.SecondaryAl[i])
		if i < last {
			body.WriteString(string(a) + `,`)
		} else {
			body.WriteString(string(a))
		}
	}
	body.WriteString(`}`)
	body.WriteString(`} `)

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       strings.NewReader(body.String()),
		Refresh:    "true",
	}

	return makeESIndexRequest(req)

}
