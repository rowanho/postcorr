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
    
    if newMax > *maxSoFar {
        *maxSoFar = newMax
        *maxI = i
        *maxJ = j
    }
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

func NeedlemanWunsch(matchReward float64, gapCost float64, a []rune, b[]rune) (float64,[]int, []int ) {
    lenA = len(a)
    lenB = len(b)
    
    d = make([][]float64, lenA +  1)
    r = make([][]int, lenB + 1)
    
    for i := range d {
        d[i] = make([]float64, lenA + 1)
        r[i] = make([]int, lenB + 1)
    }
    
    for i := 1; i < lenA + 1; i ++ {
        d[i][0] = gapCost * i
    }
    
    for j := 1; j < lenA + 1; j ++ {
        d[0][j] = gapCost * j
    }    
    
    var maxSoFar float64
    var maxI int 
    var maxJ int
    
    for i := 1; i < lenA + 1; i ++ {
        for j := 1; j < lenB + 1; j ++ {
            match := d[i-1][j-1] + matchReward;
            del := d[i-1][j] - gapCost
            ins = d[i][j-1] - gapCost
            d[i][j] = math.Max(mach, del, ins)
            updateMax(d[i][j],i,j,&maxSoFar, &maxI, &maxJ, r)
        }
    }
    
    var indicesA []int 
    var indicesB []int
    i := lenA 
    j := lenB
    
    for i > 0 && b > 0 {
        if r[i][j] == 1 {
            i -= 1
        } else if r[i][j] == 2 {
            j -= 1
        } else {
            indicesA = append([]int{i-1}, indicesA...)
            indicesB = append([]int{j-1}, indicesB...)
            i -= 1
            j -= 1
        }
    }
    
    if i > 0 {
        for k = i; k > 0; k -- {
            indicesA = append([]int{k-1}, indicesA...)
        }
    } else if j > 0 {
        for k = j; k > 0; k -- {
            indicesB = append([]int{k-1}, indicesB...)
        }
    }
    
    return maxSoFar, indicesA, indicesB
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