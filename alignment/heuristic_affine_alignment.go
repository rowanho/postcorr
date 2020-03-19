package alignment

import (
    "postCorr/flags"
    
    "math"
)


func HeuristicAffineAlignment(matchReward int, gapOp int, gapEx int, a []rune, b []rune) (int, []int, []int) {
    k := flags.ShingleSize
    bandSize := 300
    
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
    return bandedAffineDp(matchReward, gapOp, gapEx, a, b, maxBaDiff, minBaDiff)
}

func bandedAffineDp(matchReward int, gapOp int, gapEx int, a []rune , b []rune, maxBaDiff int, minBaDiff int) (int, []int, []int) {
    lenA := len(a)
    lenB := len(b)
    
    w := maxBaDiff - minBaDiff + 1
    
    d := make([][]int, lenA + 1)
    d_bc := make([][]int, lenA + 1)
    p := make([][]int, lenA + 1)
    p_bc := make([][]int, lenA + 1)
    q := make([][]int, lenA + 1)
    q_bc := make([][]int, lenA + 1)
    
    for i := 0; i < lenA + 1; i++ {
        d[i] = make([]int, w + 2)
        d_bc[i] = make([]int, w + 2)
        p[i] = make([]int, w + 2)
        p_bc[i] = make([]int, w + 2)
        q[i] = make([]int, w + 2)
        q_bc[i] = make([]int, w + 2)
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
    
    
    for j := loDiag - 1; j < hiDiag + 2; j++ {
        p[loRow][j] = math.MinInt32
        q[loRow][j] = math.MinInt32        
    }

    for i := loRow + 1; i < hiRow + 1; i ++ {
        if loDiag > 1 {
            loDiag -= 1
        }
        if i > lenB - maxBaDiff {
            hiDiag -= 1
        }
        
        d[i][loDiag - 1] = math.MinInt32
        p[i][loDiag - 1] = math.MinInt32
        q[i][loDiag - 1] = math.MinInt32
        var match int
        for j := loDiag; j < hiDiag + 1; j++ {
            if a[i-1] == b[j+i+minBaDiff-2] {
                match = matchReward
            } else {
                match = -matchReward
            }
            q[i][j] = Max(q[i][j-1] - gapEx, q[i][j-1] - gapOp)
            p[i][j] = Max(p[i-1][j+1] - gapEx, d[i-1][j+1] - gapOp)
            d[i][j] = Max(d[i-1][j] + match, Max(p[i][j], Max(q[i][j], 0)))
            
            if score < d[i][j] {
                score = d[i][j]
                maxI = i
                maxJ = j
            }
            
            if d[i][j-1] - gapOp < q[i][j-1] - gapEx {
                q_bc[i][j] = 1
            }
            
            if  d[i-1][j+1] - gapOp < p[i-1][j+1] - gapEx {
                p_bc[i][j] = 1
            }
            
            if d[i][j] == d[i-1][j] + match {
                d_bc[i][j] = 1
            }
            
            if d[i][j] == p[i][j] {
                d_bc[i][j] += 10
            }
            
            if d[i][j] == q[i][j] {
                d_bc[i][j] += 100
            }            
        }
    }
    
    aIndices := make([]int, 0)
    bIndices := make([]int, 0)
    
    i := maxI
    j := maxJ
    level := 0
    for d[i][j] > 0 {
        if level == 0 {
            if (d_bc[i][j] / 10) % 10 == 1 {
                level = 1
            } else if (d_bc[i][j] / 100) % 10 == 1 {
                level = 2
            } else if d_bc[i][j] % 10 == 1 {
                aIndices = append(aIndices, i-1)
                bIndices = append(bIndices, j + i + minBaDiff - 2)
                i -= 1
            }
        } else if level == 1 {
            if p_bc[i][j] == 0 {
                level = 0
            }
            i -= 1
            j += 1
        } else {
            if q_bc[i][j] == 0 {
                level = 0
            }
            j -= 1                
        }
    }

    return score, reverseInt(aIndices), reverseInt(bIndices)
}
