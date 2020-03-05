package correction

import (
	"postCorr/common"
	"postCorr/flags"
	
	"strings"
)

var reuseGraph = make(map[string][]map[string]string)

/**
*   Performs a majority vote across all parts of the alignment
*   If indices were counted as aligning, they are used in the vote
*   The relationship between alignments in a cluster is such that
*   the primary alignment region is very similar in both
*   Also eturns an integer representing the number of corrections made
**/

func MajorityVote(primaryDocumentID string, alignmentMaps []alignMap, documents []common.Document, docMap map[string]int) ([]rune, int) {
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
		for _, alMap := range alignmentMaps {
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
		//fmt.Println(counts)
		//fmt.Println(primText[ind])

		if primText[ind] != maxRune {
			primText[ind] = maxRune
			noCorrections += 1
		}
	}
	//fmt.Println(string(primText))
	if flags.WriteData && noCorrections > 0 {
		reuseCluster := make(map[string]string)
		p := []rune(strings.Repeat("_", maxEnd + 1 - minStart))
		for _, m := range(alignmentMaps) {
			for i := m.Start; i <= m.End; i++ {
				if _, exists := m.Mapping[i]; exists {
					p[i - minStart] = primText[i]
				}
			}
		}
		reuseCluster[primaryDocumentID] = string(p)
		for _, m := range(alignmentMaps) {
			s := m.Mapping[m.Start]
			e := m.Mapping[m.End]
			r := strings.Repeat("_", m.Start - minStart)
			t := []rune(strings.Repeat("_", e + 1 - s))
			for _, secPos := range m.Mapping {
				t[secPos - s] = documents[docMap[m.SecondaryDocumentID]].Text[secPos]
			}
			reuseCluster[m.SecondaryDocumentID] = r + string(t)
		}
		
		reuseGraph[primaryDocumentID] = append(reuseGraph[primaryDocumentID], reuseCluster)
	}
	return primText, noCorrections
}
