package alignment

import (
	"postCorr/common"

	"sync+"
)

func AlignParallel(likelyMatchingDocs map[string]map[string]bool, docs []common.Document, docMap map[string]int) map[string][]common.Alignment {

	alignments := make(map[string][]common.Alignment, 0)
	for primID, secIDs := range likelyMatchingDocs {
		primDoc := docs[docMap[primID]]
		for secID, _ := range secIDs {
			if _, exists := likelyMatchingDocs[secID][primID]; exists {
				delete(likelyMatchingDocs[secID], primID)
			}
		}

		alignmentChannel := make(chan []common.Alignment)
		for secID, _ := range secIDs {
			secDoc, _ := docs[docMap[secID]]
			go func(channel chan []common.Alignment, primDoc common.Document, secDoc common.Document, secID string) {
				alignments, inverseAlignments := alignment.GetAlignments(1.0, 2.0, primDoc, secDoc, 1, 0.0)
				channel <- alignments
				channel <- inverseAlignments
			}(alignmentChannel, primDoc, secDoc, secID)

		}
		alignments[primID] = make([]common.Alignment, 0)
		for i := 0; i < len(secIDs)*2; i++ {
			als := <-alignmentChannel
			alignments[primID] = append(alignments[primID], als...)
		}
	}
	return alignments

}

func AlignSerial(likelyMatchingDocs map[string]map[string]bool, docs []common.Document, docMap map[string]int) map[string][]common.Alignment {

	alignments := make(map[string][]common.Alignment, 0)
	for primID, secIDs := range likelyMatchingDocs {
		primDoc := docs[docMap[primID]]
		for secID, _ := range secIDs {
			if _, exists := likelyMatchingDocs[secID][primID]; exists {
				delete(likelyMatchingDocs[secID], primID)
			}
		}

		alignments[primID] = make([]common.Alignment, 0)
		alignmentChannel := make(chan []common.Alignment)
		for secID, _ := range secIDs {
			secDoc, _ := docs[docMap[secID]]
			als, inverseAls := alignment.GetAlignments(1.0, 2.0, primDoc, secDoc, 1, 0.0)
			alignments[primID] = append(alignments[primID], als...)
			alignments[primID] = append(alignments[primID], inverseAls...)
		}
	}

	return alignments

}

func GetSimilarAlignments(alignmentMap map[string][]common.Alignment) {
    
}
