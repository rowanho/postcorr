package fingerprinting

import (
	"postCorr/common"
	"postCorr/flags"
)

// Computes the jaccardindex of the two sets of fingerprints
func fpJaccardScore(fp1 map[uint64]int, fp2 map[uint64]int) float64 {

	intersection := 0
	union := 0

	// Iterate over the hashes
	for hash := range fp1 {
		if _, ok := fp2[hash]; ok {
			union += 1
			intersection += 1
		} else {
			union += 1
		}
	}

	for hash2 := range fp2 {
		if _, ok := fp1[hash2]; !ok {
			union += 1
		}
	}

	return float64(intersection) / float64(union)
}

func getSimilarLsh(docs []common.Document) map[string]map[string]bool {

	fps := make([]common.LSH_fp, len(docs))
	for i, doc := range docs {
		fp := MinHash(doc.ID, string(doc.Text), 7)
		fps[i] = fp
	}
	IndexMinHashObject()

	documentAdjacencyList := make(map[string]map[string]bool)
	for i, fp := range fps {
		sameBucketIds := SameBucketIds(fp.Signature)
		for _, id := range sameBucketIds {
			if id != docs[i].ID {
				documentAdjacencyList[docs[i].ID][id] = true
			}

		}
	}
	return documentAdjacencyList
}

func getSimilarModP(docs []common.Document) map[string]map[string]bool {
	fps := make([]map[uint64]int, len(docs))
	for i, doc := range docs {
		fp := ModP(string(doc.Text), 7, 2)
		fps[i] = fp
	}
	documentAdjacencyList := make(map[string]map[string]bool)
	for i, fp1 := range fps {
		for j, fp2 := range fps {
			if i == j {
				continue
			}
			if fpJaccardScore(fp1, fp2) > flags.JaccardThreshold {
				documentAdjacencyList[docs[i].ID][docs[j].ID] = true
			}
		}
	}
	return documentAdjacencyList
}

func GetSimilarDocuments(docs []common.Document) map[string]map[string]bool {
	var documentAdjacencyList map[string]map[string]bool
	if flags.FpType == common.MinhashFP {
		documentAdjacencyList = getSimilarLsh(docs)
	} else if flags.FpType == common.ModFP {
		documentAdjacencyList = getSimilarModP(docs)
	}

	return documentAdjacencyList
}
