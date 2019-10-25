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
        return -matchReward
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


/**
Function SmithWaterman - Computes the best local alignment for two strings, based on scoring inputs
parameter r1 - Array of runes, the first string
parameter r2  - Array of runes, the second string
parameter matchReward - The score benefit of matching two characters together
parameter gapCost - The cost of leaving a gap of length 1
returns int - The i position in the scoring matrix with the maximum score
returns int - The j position in the scoring matrix with the maximum score
returns float64 - The maximum score
returns [][]int - The scoring matrix
**/

func SmithWaterman(r1 []rune, r2 []rune, matchReward float64, gapCost float64) (int, int, float64, [][]int) {
    var cost float64


    r1Len := len(r1)
    r2Len := len(r2)
    // Initialise the scoring matrix
    d := make([][]float64, r1Len + 1)
    // Initialise the matrix for reconstructing the alignment
    a := make([][]int, r1Len + 1)
    for i := range d {
        d[i] = make([]float64, r2Len + 1)
        a[i] = make([]int, r2Len + 1)
    }

    var maxSoFar float64
    var maxI int
    var maxJ int
    

    for i := 1; i < r1Len + 1; i++ {
        for j := 1; j < r2Len + 1; j++ {
            cost = getCost(r1, i-1, r2, j-1, matchReward, gapCost)

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

func ReconstructSolution(r1 []rune, r2 []rune, maxI int, maxJ int, a[][]int) []rune{
    i := maxI 
    j := maxJ
    var solution []rune
    
    for i > 0 && j > 0 {
        if a[i][j] == 1 {
            i -= 1
        } else if a[i][j] == 2{
            j -= 1
        } else {
            solution = append( []rune{r1[i-1]}, solution ...)
            i -= 1
            j -= 1
        }
    }
    
    return solution
}