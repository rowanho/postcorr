package alignment

import (
    "postCorr/flags"
    "postCorr/fingerprinting"
    
    "math"
)

func Max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

func Min(x, y int) int {
    if x > y {
        return y
    }
    return x
}
// Gets the frequencies of k grams
func getKwords(seq []rune, k int) map[uint64][]int{
    table := make(map[uint64][]int)
    for i:=0; i < len(seq) -k; i++ {
        fp := fingerprinting.ComputeFNV64(string(seq[i: i + k]))
        table[fp] = append(table[fp], i)
    }
    return table
}


func getDiagonalSums(start int, end int, tableA map[uint64][]int, tableB map[uint64][]int) map[int]int {
    sums := make(map[int]int)
    
    for fp, indiceListA := range tableA {
        if indiceListB, exists := tableB[fp]; exists {
            for _, a := range indiceListA {
                for _, b := range indiceListB {
                    sums[a - b] ++
                }
            }
        }
    }
    return sums
}


func updateArray(i int, j int, r [][]int, match int, delete int, insert int) {
    m := Max(match, Max(Max(delete, insert), 0))
    
    if insert == m {
        r[i][j] = 2
    } else if delete == m {
        r[i][j] = 1
    } else if match == m {
        r[i][j] = 3
    } else {
        r[i][j] = 0
    }
}

func bandedDp(matchReward int, gapCost int, a []rune , b []rune, maxBaDiff int, minBaDiff int) (int, []int, []int) {
    lenA := len(a)
    lenB := len(b)
    
    w := maxBaDiff - minBaDiff + 1
    
    h := make([][]int, lenA + 1)
    r := make([][]int, lenA + 1)
    for i := 0; i < lenA + 1; i++ {
        h[i] = make([]int, w + 2)
        r[i] = make([]int, w + 2)
    }
    
    hiDiag := w
    var loDiag int
    if minBaDiff > 0 {
        loDiag = 2
    } else if maxBaDiff < 0 {
        loDiag = hiDiag + 1
    } else {
        loDiag = 2 - minBaDiff
    }
        
    loRow := Max(0, - maxBaDiff)
    hiRow := Min(lenA, lenB - minBaDiff)
    
    score :=  math.MinInt32
    
    maxI := loRow
    maxJ := loDiag - 1
    
    h[0][loDiag - 1]  = math.MinInt32
    
    for i := loRow + 1; i < hiRow + 1; i ++ {
        if loDiag > 1 {
            loDiag -= 1
        }
        if i > lenB - maxBaDiff {
            hiDiag -= 1
        }
        
        h[i][loDiag - 1] = math.MinInt32
        
        var delete, insert, match int
        for j := loDiag; j < hiDiag + 1; j++ {
            delete = h[i-1][j+1] - gapCost
            insert = h[i][j-1] - gapCost
            if a[i-1] == b[j+i+minBaDiff-2] {
                match = h[i-1][j] + matchReward
            } else {
                match = h[i-1][j] - matchReward
            }
            h[i][j] = Max(match, Max(Max(delete, insert), 0))
            
            if score < h[i][j] {
                score = h[i][j]
                maxI = i
                maxJ = j
            }
            updateArray(i, j, r, match, delete, insert)
        }
    }
    
    aIndices := make([]int, 0)
    bIndices := make([]int, 0)
    
    i := maxI
    j := maxJ
    for r[i][j] != 0 {
        if r[i][j] == 1 {
            i -= 1
            j += 1
        } else if r[i][j] == 2 {
            j -= 1
        } else {
            aIndices = append(aIndices, i-1)
            bIndices = append(bIndices, j + i + minBaDiff - 2)
            i -= 1
        }
    }
    return score, reverseInt(aIndices), reverseInt(bIndices)
}


func findPeakRegion(diagonalSums map[int]int, width int) (int, int) {
    min := 0
    max := 0
    maxVal := 0
    for width, value := range diagonalSums {
        if width < min {
            min = width
        } 
        if width > max {
            max = width
        } 
        if value > maxVal {
            maxVal = value
        }
    }
        
    if width <= max - min {
        return maxVal, (max - min) / 2
    }
    
    bestAbDiff := 0
    maxMatchSum := 0
    for i:= min; i < min + width; i++ {
        maxMatchSum += diagonalSums[i] 
    }
    currentSum := maxMatchSum
    bestAbDiff = min / 2
    for current := min + width; current <= max; current++ {
        currentSum += diagonalSums[current]
        currentSum -= diagonalSums[current]
        if currentSum > maxMatchSum {
            maxMatchSum = currentSum
            bestAbDiff = (current - width) / 2
        }
    }
    return maxMatchSum, bestAbDiff
}


func HeuristicAlignment(matchReward int, gapCost int, a []rune, b []rune) (int, []int, []int) {
    k := flags.ShingleSize
    bandSize := 20
    
    tableA := getKwords(a, k)
    tableB := getKwords(b, k)
    diagonalSums := getDiagonalSums(k - len(b), len(a) - k, tableA, tableB)
    
    maxMatchSum, bestAbDiff := findPeakRegion(diagonalSums, bandSize)
    
    if maxMatchSum == 0 {
        return 0.0, []int{}, []int{}
    }
    
    bestBaDiff := - bestAbDiff 
    maxBaDiff := Min(bestBaDiff + bandSize / 2, len(b) - 1)
    minBaDiff := Max(bestBaDiff - bandSize / 2, 1 - len(a))
    
    if maxBaDiff - minBaDiff <= 0 {
        return 0.0, []int{}, []int{}        
    }
    return bandedDp(matchReward, gapCost, a, b, maxBaDiff, minBaDiff)
}