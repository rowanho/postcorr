package alignment

import (
	"postCorr/common"
	"github.com/google/uuid"
)



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
							 
		primaryString := primary.Text
		secondaryString := secondary.Text
		
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

