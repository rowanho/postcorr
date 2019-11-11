package alignment 

import (
    "math"
)

func ReverseRune(f []rune) []rune {
    for i, j := 0, len(f)-1; i < j; i, j = i+1, j-1 {
        f[i], f[j] = f[j], f[i]
    }
    return f
}

func ReverseFloat(f []float64) []float64 {
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
            match := score[0][j - 1] + matchReward
            delete := score[0][j] - gapCost
            insert := score[1][j-1] - gapCost
            score[1][j] = math.Max(match, math.Max(delete, insert))
        }
        score[0] = score[1]
    }
    
    return score[1]
}

func swScore( matchReward float64, gapCost float64,  a []rune, b []rune) (int, int) {
    
    max_i := 0
    max_j := 0
    max_val := 0.0 
    score := make([][]float64, 2)
    
    for i := 0; i < 2; i++ {
        score[i] = make([]float64, len(b))
    }
    
    for i := 1; i < len(a); i++ {
        for j := 1; j < len(b); j++ {
            match := score[0][j - 1] + matchReward
            delete := score[0][j] - gapCost
            insert := score[1][j-1] - gapCost
            score[1][j] = math.Max(match, math.Max(delete, insert))
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

func hirschberg(matchReward float64, gapCost float64, a []rune, b []rune, offsetA int, offsetB int) (float64, []int, []int) {
    
    lenA := len(a)
    lenB := len(b)
    
    if lenA == 0 {
        l := make([]int, lenB)
        for i := 0; i < lenB; i++ {
            l[i] = i + offsetB
        }
        return 0, []int{}, l
    } else if lenB == 0 {
        l := make([]int, lenA)
        for i := 0; i < lenA; i++ {
            l[i] = i + offsetA
        }
        return 0, l, []int{}      
    } else if lenA == 1 || lenB == 1 {
        nwRes0, nwRes1, nwRes2 := NeedlemanWunsch(matchReward, gapCost, a, b)
        listA := make([]int, lenA)
        listB := make([]int, lenB)
        for i := 0; i < lenA; i++ {
            listA[i] = nwRes1[i] + offsetA
        }
        for i := 0; i < lenB; i++ {
            listB[i] = nwRes2[i] + offsetB
        }
        return nwRes0, listA, listB
    }
    
    midA := lenA / 2
    
    lastlineLeft := nwScore(matchReward, gapCost, a[0: midA + 1], b)
    revA := ReverseRune(a)
    mid2 := lenA - midA
    lastlineRight := nwScore(matchReward, gapCost, revA[0: mid2], b)
    lastlineRight = ReverseFloat(lastlineRight)
    
    max := 0.0
    maxIndice := 0
    for i := 0; i < len(lastlineLeft); i++ {
        if max < lastlineLeft[i] + lastlineRight[i]{
            max = lastlineLeft[i] + lastlineRight[i]
            maxIndice = i
        }
    }
    midB := maxIndice
    
    firstRes0, firstRes1, firstRes2 := hirschberg(matchReward, gapCost, a[0:midA], b[0:midB], offsetA, offsetB)
    secondRes0, secondRes1, secondRes2 := hirschberg(matchReward, gapCost, a[midA:], b[midB:], midA + offsetA, midB + offsetB)

    score := firstRes0 + secondRes0
    aIndices := append(firstRes1, secondRes1...)
    bIndices := append(firstRes2, secondRes2...)
    
    return score, aIndices, bIndices
    
}
