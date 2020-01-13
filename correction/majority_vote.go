package correction

import (
	"postCorr/common"
	
	//"fmt"
)

/**
*   Performs a majority vote across all parts of the alignment
*   If indices were counted as aligning, they are used in the vote
*   The relationship between alignments in a cluster is such that
*   the primary alignment region is very similar in both
*   Also eturns an integer representing the number of corrections made
**/

func MajorityVote(primaryDocumentID string, alignmentMaps []alignMap, documents []common.Document, docMap map[string]int) (string, int) {
	noCorrections := 0
	maxEnd := 0
	minStart := 100000000

	for _, alMap := range alignmentMaps {
		if alMap.Start < minStart {
			minStart = alMap.Start
		}
		if alMap.End > maxEnd {
			maxEnd = alMap.End
		}
	}

	primText := documents[docMap[primaryDocumentID]].Text
	for ind := minStart; ind < maxEnd; ind++ {
		counts := map[rune]int{}
		max := 1
		maxRune := primText[ind]
		counts[primText[ind]] = 1
	//	fmt.Println(len(alignmentMaps))
		for _, alMap := range alignmentMaps {
		//	fmt.Println(alMap.PrimaryDocumentID)
			if val, exists := alMap.Mapping[ind]; exists {
				r := documents[docMap[alMap.SecondaryDocumentID]].Text[val]
				_, ok := counts[r]
				if ok == true {
					counts[r] += 1
				} else {
					counts[r] = 1
				}

				if counts[r] > max {
					max = counts[r]
					maxRune = r
				}
			}
		}

		if primText[ind] != maxRune {
			noCorrections += 1
		}

		if primText[ind] != maxRune {
			primText[ind] = maxRune
			noCorrections += 1
		}
	}
	return string(primText), noCorrections
}
