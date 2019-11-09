package alignment 

import (
    "math"
    "sort"
)

func nwScore( a []rune, b []rune, matchReward float64, gapCost float64) []int {
    
    score := make([][]float64, 2)
    
    for i := 0; i < 2; i++ {
        score[i] = make([]float64, len(b))
    }
    
    for j := 1; j < len(b); j++ {
        score[0][j] = score[0][j-1] - gapCost
    }
    for i := 1; i < len(a); i++ {
        for j := 1; j < len(b); j++ {
            match := score[0][j - 1] + matchReward
            delete := score[0][j] - gapCost
            insert := score[1][j-1] - gapCost
            score[1][j] = math.Max(match, delete, insert)
        }
        score[0] = score[1]
    }
    
    return score[1]
}

func swScore( a []rune, b []rune, matchReward float64, gapCost float64) (int, int) {
    
    max_i := 0
    max_j := 0
    max_val := 0 
    score := make([][]float64, 2)
    
    for i := 0; i < 2; i++ {
        score[i] = make([]float64, len(b))
    }
    
    for i := 1; i < len(a); i++ {
        for j := 1; j < len(b); j++ {
            match := score[0][j - 1] + matchReward
            delete := score[0][j] - gapCost
            insert := score[1][j-1] - gapCost
            score[1][j] = math.Max(match, delete, insert)
            if score[1][j] > max_val{
                max_val = score[1][j]
                max_i = i
                max_j = j
            }
        }
        score[0] = score[1]
    }
    
    return max_i, max_j
}

func hirschberg(matchReward float64, gapCost float64, a []rune, b []rune, offsetA int, offsetB int) (int, []int, []int) {
    
    var indicesA []int 
    var indicesB []int
    maxScore := 0
    lenA := len(a)
    lenB := len(b)
    
    if lenA == 0 {
        l := make([]int, lenB)
        for i := 0; i < lenB; i++ {
            l[i] = i + offsetB
        }
        return 0, [], l
    } else if lenB == 0 {
        l := make([]int, lenA)
        for i := 0; i < lenA; i++ {
            l[i] = i + offsetA
        }
        return 0, l, []      
    } else if lenA == 1 || lenB == 1 {
        nwResult = NeedlemanWunsch(matchReward, gapCost, a, b)
        listA := make([]int, lenA)
        listB := make([]int, lenB)
        for i := 0; i < lenA; i++ {
            listA[i] = nwResult[1][i] + offsetA
        }
        for i := 0; i < lenB; i++ {
            listB[i] = nwResult[2][i] + offsetB
        }
        return nwResult[0], listA, listB
    }
    
    midA := lenA / 2
    
    lastlineLeft := nwScore(matchReward, gapCost, a[0: midA + 1], b)
    revA = sort.Reverse(a)
    mid2 := lenA - midA
    lastlineRight := nwScore(matchReward, gapCost, a[0: mid2])
    lastlineRight := sort.Reverse(a)
    
    max := 0
    maxIndice := 0
    for i = 0; i < len(lastlineLeft){
        if max < lastlineLeft[i] + lastlineRight[i]{
            max = lastlineLeft[i] + lastlineRight[i]
            maxIndice = i
        }
    }
    midB := maxIndice
    
    firstRes = hirschberg(matchReward, gapCost, a[0:midA], b[0:midB], aOffset, bOffset)
    secondRes = hirschberg(matchReward, gapCost, a[midA:], b[midB:], midA + aOffset, midB + bOffset)

    score := firstRes[0] + secondRes[0]
    aIndices := append(firstRes[1], secondRes[1]...)
    bIndices := append(firstRes[2], secondRes[2])
    
    return score, aIndices, bIndices
    
}
