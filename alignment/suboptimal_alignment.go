package alignment

import (
	"postCorr/common"
	"github.com/google/uuid"
)


/**
* Given a list of indices, that may span multiple components,
* generates a list of per-component indices instead
* We can then use this in our alignment representation
**/

func perComponentIndices(indices []int, componentLengths []int, startIndex int, endIndex int,
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
		for indice >= (lastEnd + cLen) {
			c += 1
			if i == 0 {
				firstAffected = c
			}
			lastAffected = c
			if len(currentIndices) > 0 {
				newIndices = append(newIndices, currentIndices)
				currentIndices = []int{}
			}
			lastEnd += cLen
			cLen = componentLengths[c]
		}
		currentIndices = append(currentIndices, indice-lastEnd)
	}
	if firstAffected == -1 {
		firstAffected = 0
	}
	newIndices = append(newIndices, currentIndices)
	return newIndices, startCompIndex + firstAffected, startCompIndex + lastAffected // TODO: Fix broken login with first ans last affected
}



func createAlignment(score float64, primID string, secID string, primAl []int, secAl []int) common.Alignment {

	a := common.Alignment{
		ID: 				 uuid.New().String(), 				 
		Score:               score,
		PrimaryAl:           primAl,
		PrimaryDocumentID:   primID,

		PrimaryStartIndex:     primAl[0],
		PrimaryEndIndex:       primAl[len(primAl)-1],

		SecondaryAl:           secAl,
		SecondaryDocumentID:   secID,
		
		SecondaryStartIndex: secAl[0],
		SecondaryEndIndex: secAl[len(secAl)-1],
	}
	return a
}

type Inc = struct{
	Point int
	Amount int
}


func rescoreIndices(indices []int, increments []Inc) []int {
	 newIndices := make([]int, len(indices))
	 copy(newIndices, indices)
	 
	 for _, increment := range increments {
		 for i,_ := range newIndices {
		 	if newIndices[i] >= increment.Point {
				newIndices[i] += increment.Amount
			}
	 	}
	}
	 
	 return newIndices
}


func GetAlignments(matchReward float64, gapCost float64, primary common.Document, 
				  		 secondary common.Document, stopAt int) []common.Alignment {
							 
		primaryString := primary.AllStrings()
		secondaryString := secondary.AllStrings()
		
		count := 0 
		primIncrements := []Inc{Inc{Point:0, Amount: 0,}}
		secIncrements  :=  []Inc{Inc{Point:0, Amount: 0,}}

		alignments := make([]common.Alignment, 0)
		
		for count < stopAt && len(primaryString) > 0 && len(secondaryString) > 0{
			score, primIndices, secIndices := SmithWaterman(matchReward, gapCost, primaryString, secondaryString)
			newPrimIndices := rescoreIndices(primIndices, primIncrements)
			newSecIndices := rescoreIndices(secIndices, secIncrements)
			
			// Insert into increments in a sorted manner
			n :=  Inc{Point:0, Amount: 0,}
			primIncrements = append(primIncrements, n)
			secIncrements = append(secIncrements, n)
			
			posToInsert := 0
			for i, inc := range primIncrements {
				if inc.Point > newPrimIndices[0]{
					posToInsert = i
					break;
				}
			}
			
			copy(primIncrements[posToInsert + 1:], primIncrements[posToInsert:])
			primIncrements[posToInsert] = Inc{
				Point: newPrimIndices[0],
				Amount: len(newPrimIndices),
			}
			
			posToInsert = 0
			for j, inc := range secIncrements {
				if inc.Point > newSecIndices[0]{
					posToInsert = j
					break; 
				}	
			}
			
			copy(secIncrements[posToInsert + 1:], secIncrements[posToInsert:])
			secIncrements[posToInsert] = Inc{
				Point: newSecIndices[0],
				Amount: len(newSecIndices),
			}
			count += 1			
			al := createAlignment(score, primary.ID, secondary.ID, newPrimIndices, newSecIndices)
			
			alignments = append(alignments, al)
			
			primaryString = append(primaryString[:primIndices[0]], 
								   primaryString[primIndices[len(primIndices) -1] + 1:]...)
			secondaryString = append(secondaryString[:secIndices[0]],
									 secondaryString[secIndices[len(secIndices) -1] + 1:]...)
		} 
		return alignments
}

