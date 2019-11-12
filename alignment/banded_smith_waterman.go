package alignment 

import (
    "math"
)

func reverseRune(r []rune) []rune {
    f := append([]rune{}, r...)
    for i, j := 0, len(f)-1; i < j; i, j = i+1, j-1 {
        f[i], f[j] = f[j], f[i]
    }
    return f
}

func reverseFloat(r []float64) []float64 {
    f := append([]float64{}, r...)
    for i, j := 0, len(f)-1; i < j; i, j = i+1, j-1 {
        f[i], f[j] = f[j], f[i]
    }
    return f
}

func nwScore( matchReward float64, gapCost float64, a []rune, b []rune) []float64 {
    
    score := make([][]float64, 2)
    
    for i := 0; i < 2; i++ {
        score[i] = make([]float64, len(b))
    }
    
    for j := 1; j < len(b); j++ {
        score[0][j] = score[0][j-1] - gapCost
    }
    for i := 1; i < len(a); i++ {
        for j := 1; j < len(b); j++ {
            var match float64
            if a[i] == b[j] {
                match = score[0][j - 1] + matchReward                
            } else{
                match = score[0][j - 1] - matchReward                                
            }
            delete := score[0][j] - gapCost
            insert := score[1][j-1] - gapCost
            score[1][j] = math.Max(match, math.Max(delete, insert))
        }
        for j:= 1; j < len(b); j++ {
            score[0][j] = score[1][j]
        }
    }

    return score[1]
}

func swScore( matchReward float64, gapCost float64,  a []rune, b []rune) (int, int) {
    
    max_i := 0
    max_j := 0
    max_val := 0.0 
    score := make([][]float64, 2)
    
    lenA := len(a)
    lenB := len(b)
    for i := 0; i < 2; i++ {
        score[i] = make([]float64, len(b))
    }
    
    for i := 1; i < lenA; i++ {
        for j := 1; j < lenB; j++ {
            var match float64
            if a[i] == b[j] {
                match = score[0][j - 1] + matchReward                
            } else{
                match = score[0][j - 1] -  matchReward                                
            }
            delete := score[0][j] - gapCost
            insert := score[1][j-1] - gapCost
            score[1][j] = math.Max(match, math.Max(delete, math.Max(0.0, insert)))
            if score[1][j] > max_val{
                max_val = score[1][j]
                max_i = i
                max_j = j
            }
        }
        
        for j:= 1; j < lenB; j++ {
            score[0][j] = score[1][j]
        }

    }
    
    return max_i, max_j
}

func hirschberg(matchReward float64, gapCost float64, a []rune, b []rune, offsetA int, offsetB int) (float64, []int, []int) {
    lenA := len(a)
    lenB := len(b)
    
    if lenA == 0 {
        l := make([]int, lenB)
        score := 0.0
        for i := 0; i < lenB; i++ {
            l[i] = i + offsetB
            score -= gapCost
        }
        return score, []int{}, l
    } else if lenB == 0 {
        l := make([]int, lenA)
        score := 0.0
        for i := 0; i < lenA; i++ {
            l[i] = i + offsetA
            score -= gapCost
        }
        return score, l, []int{}      
    } else if lenA <= 3 || lenB <= 3 {
        score, nwResA, nwResB := NeedlemanWunsch(matchReward, gapCost, a, b)
        listA := make([]int, len(nwResA))
        listB := make([]int, len(nwResB))
        for i := 0; i < len(nwResA); i++ {
            listA[i] = nwResA[i] + offsetA
        }
        for i := 0; i < len(nwResB); i++ {
            listB[i] = nwResB[i] + offsetB
        }
        return score, listA, listB
    }
    
    midA  := lenA / 2
    
    lastlineLeft := nwScore(matchReward, gapCost, a[0: midA], b)
    revASlice := reverseRune(a[midA:])
    revB := reverseRune(b)
    lastlineRight := nwScore(matchReward, gapCost, revASlice, revB)
    lastlineRight = reverseFloat(lastlineRight)
    
    max := 0.0
    maxIndice := 0
    for i := 0; i < len(lastlineLeft); i++ {
        if max <= lastlineLeft[i] + lastlineRight[i]{
            max = lastlineLeft[i] + lastlineRight[i]
            maxIndice = i
        }
    }
    midB := maxIndice 
    firstScore, firstRes1, firstRes2 := hirschberg(matchReward, gapCost, a[0:midA], b[0:midB], offsetA, offsetB)
    secondScore, secondRes1, secondRes2 := hirschberg(matchReward, gapCost, a[midA:], b[midB:], midA  + offsetA, midB + offsetB)

    score := firstScore + secondScore
    aIndices := append(firstRes1, secondRes1...)
    bIndices := append(firstRes2, secondRes2...)
    
    return score, aIndices, bIndices
    
}

func SmithWaterman(matchReward float64, gapCost float64, a []rune, b []rune) (float64, []int, []int){
    if len(a) == 0 || len(b) == 0 {
        return 0.0, []int {}, []int {}
    }
    endA, endB := swScore(matchReward, gapCost, a, b)
    revA := reverseRune(a)
    revB := reverseRune(b)
    revStartA, revStartB := swScore(matchReward, gapCost, revA, revB)
    startA := len(a) - 1 - revStartA    
    startB := len(b) - 1 - revStartB
    
    return hirschberg(matchReward, gapCost, b[startB: endB + 1], a[startA: endA +1], startB, startA)        
}
