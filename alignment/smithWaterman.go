package alignment


import (
    "math"
)
/**
Given 2 arrays of runes and 2 indexes, returns the cost at the respective indexes
**/

func getCost(r1 []rune, r1Index int, r2 []rune, r2Index int, matchReward float64, gapCost float64) float64 {
    if r1[r1Index] == r2[r2Index] {
        return matchReward
    } else {
        return -gapCost - 1
    }
    
}

/**
Updates the maximum value of the alignment found, and adds to array a accordingly
**/
func updateMax(newMax float64, i int, j int, maxSoFar *float64, maxI *int, maxJ *int, a [][]int) {
    
    if  *maxI < i && *maxJ == j{
        a[i][j] = 1
    } else if *maxJ < j && *maxI == i{
        a[i][j] = 2
    } else {
        a[i][j] = 3    
    }
    
    *maxSoFar = newMax
    *maxI = i
    *maxJ = j
}




func SmithWaterman(r1 []rune, r2 []rune, matchReward float64, gapCost float64) (int, int, float64, [][]int) {
    var cost float64


    r1Len := len(r1)
    r2Len := len(r2)

    // Initialise the scoring matrix
    d := make([][]float64, r1Len)
    // Initialise the matrix for reconstructing the alignment
    a := make([][]int, r1Len)
    for i := range d {
        d[i] = make([]float64, r2Len)
        a[i] = make([]int, r1Len)
    }

    var maxSoFar float64
    var maxI int
    var maxJ int
    
    for i := 0; i < r1Len; i++ {
        // substitution cost
        cost = getCost(r1, i, r2, 0, matchReward, gapCost)
        if i == 0 {
            d[0][0] = math.Max(0.0, math.Max(-gapCost, cost))
        } else {
            d[i][0] = math.Max(0.0, math.Max(d[i-1][0]-gapCost, cost))
        }

        // save if it is the biggest thus far
        if d[i][0] > maxSoFar {
            updateMax(d[i][0],i,0,&maxSoFar, &maxI, &maxJ, a)
        }
    }

    for j := 0; j < r2Len; j++ {
        // substitution cost
        cost = getCost(r1, 0, r2, j, matchReward, gapCost)
        if j == 0 {
            d[0][0] = math.Max(0, math.Max(-gapCost, cost))
        } else {
            d[0][j] = math.Max(0, math.Max(d[0][j-1]-gapCost, cost))
        }

        // save if it is the biggest thus far
        if d[0][j] > maxSoFar {
            updateMax(d[0][j],0,j,&maxSoFar, &maxI, &maxJ, a)
        }
    }

    for i := 1; i < r1Len; i++ {
        for j := 1; j < r2Len; j++ {
            cost = getCost(r1, i, r2, j, matchReward, gapCost)

            // find the lowest cost
            d[i][j] = math.Max(
                math.Max(0, d[i-1][j]-gapCost),
                math.Max(d[i][j-1]-gapCost, d[i-1][j-1]+cost))

            // save if it is the biggest thus far
            if d[i][j] > maxSoFar {
                updateMax(d[i][j],i,j,&maxSoFar, &maxI, &maxJ, a)
            }
        }
    }

    return maxI, maxJ, maxSoFar, a
}

func ReconstructSolution(r1 []rune, r2 []rune, maxI int, maxJ int, a[][]float64) []rune{
    i := maxI 
    j := maxJ
    var solution []rune
    
    for i > 0 && j > 0 {
        if a[i][j] == 1 {
            i -= 1
        } else if a[i][j] == 2{
            j -= 1
        } else {
            solution = append(solution, r1[i])
            i -= 1
            j -= 1
        }
    }
    
    return solution
}