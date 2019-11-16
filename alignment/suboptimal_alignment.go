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

func contains(s []int, e int) int {
    for i, a := range s {
        if a == i {
            return i
        }
    }
    return - 1
}


// Fills the missing text of a segment
func fillSegmentText(s Segment, componentOrder []string, components map[string][]rune) {
	
	if s.StartComp == s.EndComp {
		comp := components[componentOrder[s.StartComp]]
		s.Text = comp[s.StartPos: s.EndPos + 1]
	}
	for i := s.StartComp; i < s.EndComp + 1; i ++ {
		comp := components[componentOrder[i]]
		if i == s.StartComp {
			s.Text = append(s.Text, comp[s.StartPos:]...)			
		} else if i == s.EndComp {
			s.Text = append(s.Text, comp[:s.EndPos + 1]...)
		} else {
			s.Text = append(s.Text, comp...)
		}
	}
}

/**
* Gets a new list of segmnets after we remove the segment with the alignment
* The removed segment has to be a sub-segment of one of the previous list
**/

func getNewSegments(prevSegments []Segment, affectedSegs []int, removedSegs []Segment, componentOrder []string, components map[string][]rune,
					) []Segment {
	
	segs := make([]Segment, 0)
	for i, s := range prevSegments {
		// If affected
		if contains(affectedSegs, i) != -1 {
			r := removedSegs[i]
			// Build the before segment
			var before Segment
			before.StartComp = s.StartComp
			before.StartPos = s.StartPos
			if s.StartComp < r.StartComp {
				if r.EndPos > 0 {
					before.EndComp = r.StartComp
					before.EndPos = r.StartPos -1
				} else {
					before.EndComp = r.StartComp - 1
					l := len(components[componentOrder[before.EndComp]])
					before.EndPos = l - 1
				}
				fillSegmentText(before, componentOrder, components)
				segs = append(segs, before)
				
			} else {
				before.EndComp = s.StartComp
				if s.StartPos < r.StartPos {
					before.EndPos = r.StartPos - 1
					fillSegmentText(before, componentOrder, components)
					segs = append(segs, before)
				}
			}
			
			// Build the 'after' segment
			var after Segment
			
			
			after.EndPos = s.EndPos
			after.EndComp = s.EndComp					
			if r.EndComp < s.EndComp {
				if r.EndPos == len(components[componentOrder[r.EndComp]]) - 1{
					after.StartComp = r.EndComp + 1
					after.StartPos = 0
				} else {
					after.StartComp = r.EndComp
					after.StartPos = r.EndPos + 1
				}
				fillSegmentText(after, componentOrder, components)
				segs = append(segs, after)
	
			} else {
				after.StartComp = s.EndComp
				if r.StartPos < s.StartPos {
					before.StartPos = r.EndPos + 1
					fillSegmentText(after, componentOrder, components)
					segs = append(segs, after)				
				}
			}
		} else {
			segs = append(segs, s)
		}
	 
	}
	return segs
}


/**
* Given a list of indices, that may span multiple components,
* generates a list of per-component indices instead
* We can then use this in our alignment representation
**/

func rescoreIndices(indices []int, componentLengths []int, startIndex int, endIndex int, 
					startCompIndex int) ([][]int, int, int) {
	newIndices := make([][]int, 0)
	
	lastEnd := 0
	
	currentIndices := make([]int, 0)
	c := 0
	cLen := componentLengths[c]
	
	firstAffected := -1
	lastAffected := 0
	for i, indice := range indices {
		// Now on the next component
		for indice >= (lastEnd + cLen){
			c += 1
			if i == 0{
				firstAffected = c
			} 
			lastAffected = c
			if len(currentIndices) > 0{
				newIndices = append(newIndices, currentIndices)
				currentIndices = []int{}
			}
			lastEnd += cLen
			cLen = componentLengths[c]
		}
		currentIndices = append(currentIndices, indice - lastEnd)
	}
	if firstAffected == -1 {
		firstAffected = 0
	}
	newIndices = append(newIndices, currentIndices)	
	return newIndices, startCompIndex + firstAffected, startCompIndex + lastAffected // TODO: Fix broken login with first ans last affected
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
			EndComp : len(doc.ComponentOrder) -1,
			StartPos : 0,
			EndPos : lastCompLen - 1,
			Text : doc.AllStrings(),
		},
	}
	return initialSegments
}


func getComponentLengths(doc common.Document) []int {
	lengths := make([]int, len(doc.ComponentOrder))
	
	for i, compID :=  range doc.ComponentOrder {
		lengths[i] = len(doc.TextComponents[compID])
	}
	return lengths
}



func createAlignment(score float64, primID string, secID string, primAls [][]int, secAls [][]int,
	 				primCompIDs []string, secCompIDs []string) common.Alignment {
	
	lastAlign := primAls[len(primAls) - 1]
	
	a := common.Alignment{
		Score: score,
		PrimaryAl : primAls,
		PrimaryDocumentID : primID,
		PrimaryComponentIDs : primCompIDs,
		
		PrimaryStartComponent : primCompIDs[0],
		PrimaryEndComponent : primCompIDs[len(primCompIDs) -1],
		PrimaryStartIndex : primAls[0][0],
		PrimaryEndIndex : lastAlign[len(lastAlign) -1],
		
		SecondaryAl: secAls,
		SecondaryDocumentID: secID,
		SecondaryComponentIDs: secCompIDs,
		
	}
	return a
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
	
	for count < stopAt  {
		
		affectedPrim := make([]int, 0)
		affectedSec := make([]int, 0)
		
		removedPrimSegs := make([]Segment, 0)
		removedSecSegs := make([]Segment, 0)
		
		for i, pSeg := range primarySegments {
			for j, sSeg := range secondarySegments {
				score, primIndices, secIndices := SmithWaterman(matchReward, gapCost, pSeg.Text, sSeg.Text)
				
				if len(primIndices) == 0 {
					continue;
				}
				rescoredPrims, primFirstAffected, primLastAffected := rescoreIndices(primIndices, primCompLengths[pSeg.StartComp:pSeg.EndComp + 1 ],
					 																 pSeg.StartPos, pSeg.EndPos, pSeg.StartComp)
				rescoredSecs, secFirstAffected, secLastAffected := rescoreIndices(secIndices, secCompLengths[sSeg.StartComp:sSeg.EndComp + 1 ],
					 															  sSeg.StartPos, sSeg.EndPos, sSeg.StartComp)

				al := createAlignment(score, primary.ID, secondary.ID, rescoredPrims, rescoredSecs, 
									  primary.ComponentOrder[pSeg.StartComp:pSeg.EndComp + 1 ],
								  	  secondary.ComponentOrder[sSeg.StartComp:sSeg.EndComp + 1 ])
				
				
				lastPrim := rescoredPrims[len(rescoredPrims) -1]
				rPrim := Segment {
					StartComp : primFirstAffected, 
					EndComp : primLastAffected,
					StartPos : rescoredPrims[0][0],
					EndPos : lastPrim[len(lastPrim) -1],
				}
				
				lastPrim = rescoredSecs[len(rescoredSecs) -1]
				rSec := Segment {
					StartComp : secFirstAffected, 
					EndComp : secLastAffected,
					StartPos : rescoredSecs[0][0],
					EndPos : lastPrim[len(lastPrim) -1],
				}
				
				removedPrimSegs = append(removedPrimSegs, rPrim)
				removedSecSegs = append(removedSecSegs, rSec)
								  
				alignments = append(alignments, al)
				affectedPrim = append(affectedPrim, i)
				affectedSec = append(affectedSec, j)

			}
		}
		primarySegments = getNewSegments(primarySegments, affectedPrim, removedPrimSegs, primary.ComponentOrder, primary.TextComponents)
		secondarySegments = getNewSegments(secondarySegments, affectedSec, removedSecSegs, secondary.ComponentOrder, secondary.TextComponents)
		affectedPrim = make([]int, 0)
		affectedSec = make([]int, 0)
		
		removedPrimSegs = make([]Segment, 0)
		removedSecSegs = make([]Segment, 0)
		count += 1 
	}
	return alignments

}

