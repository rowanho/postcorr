package alignment

import (
	"postCorr/common"
)

// We use segments to track continuous groups of text
// They have  a start and end component indice
// As well as a start and end position

type Segment = struct {
	StartComp int
	EndComp int
	StartPos int
	EndPos int
	Text []rune
}

/**
* Gets a new list of segmnets after we remove the segment with the alignment
* The removed segment has to be a sub-segment of one of the previous list
**/

func getNewSegments(prevSegments []Segment, removedSeg Segment, componentOrder []string, components map[string][]rune
					) []Segment {
	
	segs := make([]Segment, 0)
	for _, s := range prevSegments {
		// If removedSeg is a sub-segment of r
		if r.StartComp <= removedSeg.StartComp && r.EndComp >= removedSeg.EndComp{
			if r.StartPos <= removedSeg.StartPos && r.EndPos >= removedSeg.EndPos {
				// Get the other two sub segments, if they are big enough
				
			} else{
				segs = append(segs, r)
			}
			
		} else {
			segs = append(segs, r)
		} 
	}
	return segs
}


/**
* Given a list of indices, that may span multiple components,
* generates a list of per-component indices instead
* We can then use this in our alignment representation
**/

func rescoreIndices(indices []int, componentLengths []int, startIndex int, endIndex int 
					) [][]int {
	newIndices := make([][]int, 0)
	
	lastEnd := 0
	
	currentIndices := make([]int)
	c := 0
	cLen := componentLengths[c]
	for _, indice := range indices {
		// Now on the next component
		if indice > lastEnd + cLen - 1{
			c += 1
			newIndices := append(newIndices, currentIndices)
			currentIndices = []int{}
			lastEnd += cLen
			cLen = componentLengths[c]
		}
		currentIndices = append(currentIndices, indice - lastEnd)
	}
	
	return newIndices
}

/**
* The initial segment of a document, ie all strings
* Wrapped in a slice
**/

func getInitialSegment(doc common.Document) []Segment {
	lastCompID :=  doc.ComponentOrder[len(doc.ComponentOrder) -1]
	lastCompLen := len(doc.TextComponents[lastCompID])
	initialSegments:= []Segment{ Segment{
			StartComp : 0,
			EndComp : len(doc.ComponentOrder),
			StartPos : 0,
			EndPos : lastCompLen,
			Text : doc.AllStrings()
		}
	}
	
}


func getComponentLengths(doc common.Document) []int {
	lengths := make([]int, len(doc.ComponentOrder))
	
	for i, compID :=  range doc.ComponentOrder {
		lengths[i] = len(doc.TextComponents[compID])
	}
	return lengths
}



func createAlignment(primID string, secID string, primAls [][]int, secAls [][]int,
	 				primCompIDs []string, secCompIDs []string) common.Alignment {
	
	
}
/**
* Takes in two document objects
* Returns a list of the alignment objects between them
**/

func GetAlignments(matchReward float64, gapCost float64, primary common.Document, secondary common.Document, stopAt int,
				   ) []common.Alignment {
					   
	alignments := make([]common.Alignment, 0)
	
	
	primarySegments := getInitialSegment(primary)
	secondarySegments := getInitialSegment(secondary)
	
	primCompLengths := getComponentLengths(primary)
	secCompLengths := getComponentLengths(secondary)
	
	count := 0
	
	var removedPrimSeg Segment
	var removedSecSeg Segment 
	for count < stopAt  {
		for i, pSeg := range primarySegments {
			for j, sSeg := range secondarySegments {
				score, primIndices, secIndices := SmithWaterman(matchReward, gapCost, primarySegments, secondarySegments)
				rescoredPrims := rescoreIndices(indices, primCompLengths[pSeg.StartComp:pSeg.EndComp + 1 ], pSeg.StartPos, pSeg.EndPos)
				rescoredSecs := rescoreIndices(indices, secCompLengths[sSeg.StartComp:sSeg.EndComp + 1 ], sSeg.StartPos, sSeg.EndPos)
				al := createAlignment(primary.ID, secondary.ID, rescoredPrims, rescoredSecs, 
									  primary.ComponentOrder[pSeg.StartComp:pSeg.EndComp + 1 ],
								  	  secondary.ComponentOrder[sSeg.StartComp:sSeg.EndComp + 1 ])
								  
				alignments = append(alignments, al)
				
			}
		}
		primarySegments = getNewSegments(primarySegments)
		secondarySegments = getNewSegments(secondarySegments)
		count += 1 
	}
	return alignments

}

