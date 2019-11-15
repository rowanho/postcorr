package alignment

import (
	"postCorr/common"
)


/**
* Takes in two document objects
* Returns a list of the alignment objects between them
**/

func GetAlignments(matchReward float64, gapCost float64, primary common.Document, secondary common.Document, stopAt int,
				   ) []common.Alignment {
					   
	alignments := make([]common.Alignment, 0)
	
	// Stores the ranges of components that are still 'in'
	// Stores a start and end point for a component
		
	primRanges, primCompLengths := getComponentLengths(primary.ComponentOrder, primary.TextComponents)
	secRanges, secCompLengths := getComponentLengths(secondary.ComponentOrder, secondary.TextComponents)
	
	// Nested for loop, then repeat the process
	// Technically this is an n^2 operation, however n should stay small
	
	for i, pci := range primCompIndices {
		primRunes, primStringIndices, primCompIndices := allStringsInRange(primary.ComponentOrder, 
			primary.TextComponents, primCompLengths, primRanges ) 
		secRunes, secStringIndices, secCompIndices := allStringsInRange(secondary.ComponentOrder, 
			secondary.TextComponents, secCompLengths, secRanges ) 
				
		psi := primStringIndices[i]
		pr := primRunes[i]
		for j, sci := range secCompIndices {
			ssi := secStringIndices[j]
			sr := secRunes[j]
			
			score, pIndices, sIndices := SmithWaterman(matchReward, gapCost, pr, sr)
			
			pIndiceList := rescoreIndices(pIndices, primStringIndices, primCompIndices)
			sIndiceList := rescoreIndices(sIndices, secStringIndices, secCompIndices)
		
		}
	}
}

/**
* Returns two maps, which take component ids as keys
* The first is simply the lengths of the component texts
* The second is the start and end indexes, for convenience as this is needed in GetAlignments 
**/
func getComponentLengths(componentOrder []string, components map[string][]rune) (map[string][]int,map[string]int) {
	lengths := make(map[string]int)
	ranges := make(map[string][]int)
	
	for _, compID := range componentOrder {
		l := len(components[compID])
		lengths[compID] = l
		ranges[compID] = []int{0, l - 1}
	}
	
	return ranges, lengths
}


/**
* Takes in the component ids, the component lengths and their ranges
* If two consecutive components have been separated by a previous alignment, they are now considered disjoint
* Otherwise they are considered connected
* We look through connected components and build a list of the text segments we need for the actual alignment
* Each can be represented by a starting and ending component ids (stored as their indices in the list),
* As well as a starting and ending indice in that component
**/
func allStringsInRange(componentOrder []string, components map[string][]rune, lengths map[string]int, ranges map[string][]int,
					   )  ([][]rune, [][]int, [][]int) { 
	
	componentListIndexes := make([][]int, 0)
	runeIndexes := make([][]int, 0)
	runes := make([][]rune, 0)
	
	window := make([]int, 2)
	windowIndexes := make([]int, 2)
	runeBuf := make([]rune, 0)
	
	connectedPrev := false
	
	noComponents := len(componentOrder)
	
	for i, compID := range componentOrder {
			
		
		r := ranges[compID]
		l := lengths[compID]
		
		currentRune := components[compID]
		
		if r[0] > 0 {
			connectedPrev = false
		}
		
		// if we could be connected to the next component
		if (r[1] == l -1) && (i < noComponents -1) {
			if connectedPrev == true {
				runeBuf = append(runeBuf, currentRune[r[0]:]...)
				continue;
			} else {
				runeBuf = currentRune[r[0]:]
				window[0] = r[0]
				windowIndexes[0] = i
				connectedPrev = true
			}
		} else { // Not connected to the next
			if connectedPrev == true {
				runeBuf = append(runeBuf, currentRune[r[0]: r[1]]...)
				window[1] = r[1]
				windowIndexes[1] = i				
			} else {
				runeBuf = currentRune[r[0]: r[1]]
				window[0] = r[0]
				window[1] = r[1]
				windowIndexes[0] = i				
				windowIndexes[1] = i				
			}
			runeIndexes = append(runeIndexes, window)
			componentListIndexes = append(componentListIndexes, windowIndexes)
			runes = append(runes, runeBuf)
		}
	}
	
	return runes, runeIndexes, componentListIndexes
	
}

