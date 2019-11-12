// Performs suboptimal alignment between the two strings
// This finds disjoint regions of similar substrings

package alignment

func SuboptimalAlignment(matchReward float64, gapCost float64, a []rune, b []rune, stopAt int ) (scores []float64, listIndicesA [][]int, listIndicesB [][]int) {
    
    scores = make([]float64, 0)
    listIndicesA = make([][]int, 0)
    listIndicesB = make([][]int, 0)

    for i := 0; i < stopAt && len(a) > 0 && len(b) > 0; i++ {
        score,indicesA, indicesB := SmithWaterman(matchReward, gapCost, a, b)
        scores = append(scores, score)
        listIndicesA = append(listIndicesA, indicesA)
        listIndicesB = append(listIndicesB, indicesB)
        
        startA := indicesA[0]
        endA := indicesA[len(indicesA) -1]
        a = append(a[:startA + 1], a[:endA]...) 
        
        startB := indicesB[0]
        endB := indicesB[len(indicesB) -1]
        b = append(b[:startB + 1], a[:endB]...) 
    }
    
    return scores, listIndicesA, listIndicesB
    
    
}